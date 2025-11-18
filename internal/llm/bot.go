package llm

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/riz/auto-lmk/internal/model"
)

// Bot handles car sales conversations
type Bot struct {
	provider         Provider
	convRepo         ConversationRepository
	carRepo          CarRepository
	pendingImages    map[string][]string // senderPhone -> image paths
	pendingCarID     map[string]int      // senderPhone -> car_id for image context
	currentSender    string              // temporary storage for current sender during processing
}

// Import internal model types for simpler interfaces
type (
	Conversation = struct{ ID int }
	BotMessage   = struct {
		SenderPhone string
		MessageText string
		Direction   string
	}
)

// ConversationRepository interface for conversation operations
type ConversationRepository interface {
	GetOrCreate(ctx context.Context, senderPhone string, isSales bool) (*Conversation, error)
	AddMessage(ctx context.Context, conversationID int, senderPhone, messageText, direction string) error
	GetMessages(ctx context.Context, conversationID int, limit int) ([]*BotMessage, error)
}

// CarRepository interface for car operations
type CarRepository interface {
	SearchCarsForBot(ctx context.Context, filters map[string]interface{}) (interface{}, error)
	GetCarWithDetails(ctx context.Context, carID int) (interface{}, error)
	GetCarPhotos(ctx context.Context, carID int) (interface{}, error)
	CreateWithPhotos(ctx context.Context, car *model.Car, photoURLs []string) (int, error)
}

// NewBot creates a new conversation bot
func NewBot(provider Provider, convRepo ConversationRepository, carRepo CarRepository) *Bot {
	return &Bot{
		provider:      provider,
		convRepo:      convRepo,
		carRepo:       carRepo,
		pendingImages: make(map[string][]string),
		pendingCarID:  make(map[string]int),
	}
}

// ProcessMessage processes incoming message and returns bot response
func (b *Bot) ProcessMessage(ctx context.Context, tenantID int, senderPhone, messageText string, isSales bool) (string, error) {
	slog.Info("processing message", "tenant_id", tenantID, "sender", senderPhone, "is_sales", isSales)

	// Store current sender for function execution context
	b.currentSender = senderPhone

	// 1. Get conversation to access history
	conv, err := b.convRepo.GetOrCreate(ctx, senderPhone, isSales)
	if err != nil {
		slog.Error("failed to get conversation", "error", err)
		// Continue without history
	}

	// 2. Build messages for LLM
	messages := []Message{
		{
			Role:    "system",
			Content: b.buildSystemPrompt(isSales),
		},
	}

	// 3. Load recent conversation history (last 10 messages)
	if conv != nil {
		history, err := b.convRepo.GetMessages(ctx, conv.ID, 10)
		if err != nil {
			slog.Error("failed to load history", "error", err)
		} else {
			for _, msg := range history {
				role := "user"
				if msg.Direction == "outbound" {
					role = "assistant"
				}
				messages = append(messages, Message{
					Role:    role,
					Content: msg.MessageText,
				})
			}
		}
	}

	// 4. Add current user message
	messages = append(messages, Message{
		Role:    "user",
		Content: messageText,
	})

	// 5. Call LLM with functions
	functions := b.GetAvailableFunctions(isSales)
	response, err := b.provider.Chat(ctx, messages, functions)
	if err != nil {
		return "", fmt.Errorf("LLM call failed: %w", err)
	}

	// 6. Handle function calls
	if response.FunctionCall != nil {
		slog.Info("function call requested", "function", response.FunctionCall.Name, "args", response.FunctionCall.Arguments)

		// Execute function
		result, err := b.executeFunction(ctx, response.FunctionCall.Name, response.FunctionCall.Arguments)
		if err != nil {
			slog.Error("function execution failed", "error", err)
			return "Maaf, terjadi kesalahan saat memproses permintaan Anda.", nil
		}

		// Format result for LLM
		resultText := fmt.Sprintf("Hasil fungsi %s: %v", response.FunctionCall.Name, result)

		// Call LLM again with function result
		messages = append(messages, Message{
			Role:    "assistant",
			Content: response.Content,
		})
		messages = append(messages, Message{
			Role:    "function",
			Content: resultText,
		})

		response, err = b.provider.Chat(ctx, messages, functions)
		if err != nil {
			return "", fmt.Errorf("LLM call after function failed: %w", err)
		}
	}

	return response.Content, nil
}

// executeFunction executes a function called by the LLM
func (b *Bot) executeFunction(ctx context.Context, functionName string, arguments map[string]interface{}) (interface{}, error) {
	switch functionName {
	case "searchCars":
		return b.carRepo.SearchCarsForBot(ctx, arguments)

	case "getCarDetails":
		carID, ok := arguments["car_id"].(float64)
		if !ok {
			return nil, fmt.Errorf("invalid car_id")
		}
		return b.carRepo.GetCarWithDetails(ctx, int(carID))

	case "sendCarImages":
		carID, ok := arguments["car_id"].(float64)
		if !ok {
			return nil, fmt.Errorf("invalid car_id")
		}

		// Get car photos
		photos, err := b.carRepo.GetCarPhotos(ctx, int(carID))
		if err != nil {
			return nil, fmt.Errorf("failed to get car photos: %w", err)
		}

		// Store photos for sending (WhatsAppService will send them)
		if photoList, ok := photos.([]*model.CarPhoto); ok && len(photoList) > 0 {
			imagePaths := make([]string, len(photoList))
			for i, photo := range photoList {
				imagePaths[i] = photo.FilePath
			}
			b.pendingImages[b.currentSender] = imagePaths
			b.pendingCarID[b.currentSender] = int(carID)

			return map[string]interface{}{
				"status":      "images_queued",
				"image_count": len(photoList),
				"message":     fmt.Sprintf("Mengirimkan %d foto mobil...", len(photoList)),
			}, nil
		}

		return map[string]interface{}{
			"status":  "no_images",
			"message": "Maaf, foto mobil ini belum tersedia.",
		}, nil

	case "uploadCar":
		return b.executeUploadCar(ctx, arguments)

	default:
		return nil, fmt.Errorf("unknown function: %s", functionName)
	}
}

// executeUploadCar handles car upload function call from LLM
func (b *Bot) executeUploadCar(ctx context.Context, arguments map[string]interface{}) (interface{}, error) {
	// Extract and validate parameters
	brand, _ := arguments["brand"].(string)
	modelName, _ := arguments["model"].(string)
	year, _ := arguments["year"].(float64)
	price, _ := arguments["price"].(float64)
	transmission, _ := arguments["transmission"].(string)
	fuelType, _ := arguments["fuel_type"].(string)
	description, _ := arguments["description"].(string)

	// Validate required fields
	if brand == "" || modelName == "" || year == 0 || price == 0 || transmission == "" || fuelType == "" {
		return map[string]interface{}{
			"success": false,
			"error":   "Data tidak lengkap. Brand, model, tahun, harga, transmisi, dan bahan bakar wajib diisi.",
		}, nil
	}

	// Validate year range
	if year < 1990 || year > 2025 {
		return map[string]interface{}{
			"success": false,
			"error":   "Tahun harus antara 1990 dan 2025",
		}, nil
	}

	// Validate price range
	if price < 10000000 || price > 10000000000 {
		return map[string]interface{}{
			"success": false,
			"error":   "Harga harus antara 10 juta dan 10 miliar",
		}, nil
	}

	// Check for uploaded photos
	photoURLs := b.GetPendingPhotos(b.currentSender)
	if len(photoURLs) == 0 {
		return map[string]interface{}{
			"success": false,
			"error":   "Silakan upload foto mobil terlebih dahulu",
		}, nil
	}

	// Create car with photos
	car := &model.Car{
		Brand:        brand,
		Model:        modelName,
		Year:         int(year),
		Price:        int64(price),
		Transmission: &transmission,
		FuelType:     &fuelType,
		Description:  &description,
		Status:       "available",
		IsFeatured:   false,
	}

	carID, err := b.carRepo.CreateWithPhotos(ctx, car, photoURLs)
	if err != nil {
		slog.Error("failed to create car", "error", err)
		return map[string]interface{}{
			"success": false,
			"error":   fmt.Sprintf("Gagal menyimpan mobil: %v", err),
		}, nil
	}

	// Clear pending photos after successful upload
	b.ClearPendingPhotos(b.currentSender)

	// Get tenant ID for catalog URL
	tenantID, _ := model.GetTenantID(ctx)

	// Return success response
	return map[string]interface{}{
		"success":     true,
		"car_id":      carID,
		"brand":       brand,
		"model":       modelName,
		"year":        int(year),
		"price":       int(price),
		"transmission": transmission,
		"fuel_type":   fuelType,
		"photo_count": len(photoURLs),
		"catalog_url": fmt.Sprintf("https://tenant%d.auto-lmk.com/cars/%d", tenantID, carID),
		"message":     "Mobil berhasil ditambahkan ke catalog!",
	}, nil
}

// buildSystemPrompt creates system prompt based on user type
func (b *Bot) buildSystemPrompt(isSales bool) string {
	basePrompt := `Anda adalah asisten penjualan mobil yang ramah dan profesional di Auto LMK.`

	if isSales {
		return basePrompt + `

Anda sedang membantu SALES TEAM. Anda memiliki kemampuan tambahan:
- Upload mobil baru ke catalog
- Lihat dan kelola inventory

Saat sales ingin upload mobil, minta mereka:
1. Upload foto mobil (bisa multiple, maksimal 5 foto)
2. Ketik detail: "Brand Model Tahun Harga Transmisi BahanBakar"
   Contoh: "Toyota Avanza 2020 185juta AT Bensin"

Cara upload:
- Sales upload foto dulu
- Setelah foto diterima, sales ketik detail mobil
- Anda akan parse detail dan konfirmasi sebelum menyimpan
- Format harga: "185juta" atau "185jt" → 185000000
- Transmisi: "matic"/"AT" → AT, "manual"/"MT" → MT
- Bahan bakar: "bensin" → Bensin, "diesel" → Diesel

Gunakan bahasa Indonesia yang profesional dan efisien.`
	}

	return basePrompt + `

Anda membantu CUSTOMER mencari mobil. Anda dapat:
- Mencari mobil berdasarkan brand, budget, transmisi, dll
- Menampilkan detail dan foto mobil
- Memberikan rekomendasi

Jika customer tanya tentang upload atau tambah mobil, jelaskan bahwa fitur itu untuk sales team.

Gunakan bahasa Indonesia yang ramah, natural, dan helpful.
Pahami istilah automotive Indonesia seperti: matic (automatic), bensin (gasoline), OTR (On The Road price).`
}

// GetAvailableFunctions returns functions available to LLM
func (b *Bot) GetAvailableFunctions(isSales bool) []Function {
	baseFunctions := []Function{
		{
			Name:        "searchCars",
			Description: "Cari mobil yang tersedia berdasarkan kriteria",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"brand": map[string]interface{}{
						"type":        "string",
						"description": "Merek mobil (Toyota, Honda, Mitsubishi, dll)",
					},
					"model": map[string]interface{}{
						"type":        "string",
						"description": "Model mobil (Avanza, CR-V, Xpander, dll)",
					},
					"max_price": map[string]interface{}{
						"type":        "integer",
						"description": "Harga maksimal dalam Rupiah",
					},
					"transmission": map[string]interface{}{
						"type":        "string",
						"enum":        []string{"manual", "automatic"},
						"description": "Jenis transmisi (matic/automatic atau manual)",
					},
					"fuel_type": map[string]interface{}{
						"type":        "string",
						"enum":        []string{"bensin", "diesel", "hybrid", "electric"},
						"description": "Jenis bahan bakar",
					},
				},
			},
		},
		{
			Name:        "getCarDetails",
			Description: "Dapatkan detail lengkap sebuah mobil termasuk foto",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"car_id": map[string]interface{}{
						"type":        "integer",
						"description": "ID mobil",
					},
				},
				"required": []string{"car_id"},
			},
		},
		{
			Name:        "sendCarImages",
			Description: "Kirim foto-foto mobil ke customer via WhatsApp",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"car_id": map[string]interface{}{
						"type":        "integer",
						"description": "ID mobil yang fotonya akan dikirim",
					},
				},
				"required": []string{"car_id"},
			},
		},
	}

	// Add uploadCar function for sales only
	if isSales {
		baseFunctions = append(baseFunctions, Function{
			Name:        "uploadCar",
			Description: "Upload mobil baru ke catalog dengan foto dan detail",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"brand": map[string]interface{}{
						"type":        "string",
						"description": "Brand mobil (Toyota, Honda, dll)",
					},
					"model": map[string]interface{}{
						"type":        "string",
						"description": "Model mobil (Avanza, Civic, dll)",
					},
					"year": map[string]interface{}{
						"type":        "integer",
						"description": "Tahun produksi",
					},
					"price": map[string]interface{}{
						"type":        "integer",
						"description": "Harga dalam rupiah",
					},
					"transmission": map[string]interface{}{
						"type":        "string",
						"enum":        []string{"AT", "MT"},
						"description": "AT atau MT",
					},
					"fuel_type": map[string]interface{}{
						"type":        "string",
						"enum":        []string{"Bensin", "Diesel"},
						"description": "Bensin atau Diesel",
					},
					"description": map[string]interface{}{
						"type":        "string",
						"description": "Deskripsi tambahan (optional)",
					},
				},
				"required": []string{"brand", "model", "year", "price", "transmission", "fuel_type"},
			},
		})
	}

	return baseFunctions
}

// GetPendingImages returns images queued for sending to a sender
func (b *Bot) GetPendingImages(senderPhone string) []string {
	return b.pendingImages[senderPhone]
}

// ClearPendingImages clears pending images for a sender
func (b *Bot) ClearPendingImages(senderPhone string) {
	delete(b.pendingImages, senderPhone)
	delete(b.pendingCarID, senderPhone)
}

// AddPendingPhoto adds a photo to pending uploads for a sender
func (b *Bot) AddPendingPhoto(senderPhone, photoPath string) int {
	if b.pendingImages[senderPhone] == nil {
		b.pendingImages[senderPhone] = []string{}
	}
	b.pendingImages[senderPhone] = append(b.pendingImages[senderPhone], photoPath)
	return len(b.pendingImages[senderPhone])
}

// GetPendingPhotos returns pending photo uploads for a sender
func (b *Bot) GetPendingPhotos(senderPhone string) []string {
	return b.pendingImages[senderPhone]
}

// ClearPendingPhotos clears pending photo uploads for a sender
func (b *Bot) ClearPendingPhotos(senderPhone string) {
	delete(b.pendingImages, senderPhone)
}
