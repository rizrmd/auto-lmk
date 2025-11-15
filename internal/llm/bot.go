package llm

import (
	"context"
	"fmt"
	"log/slog"
)

// Bot handles car sales conversations
type Bot struct {
	provider Provider
	convRepo ConversationRepository
	carRepo  CarRepository
	leadRepo LeadRepository
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
}

// LeadRepository interface for lead operations
type LeadRepository interface {
	Create(ctx context.Context, phoneNumber, name string, interestedCarID *int) (int, error)
}

// NewBot creates a new conversation bot
func NewBot(provider Provider, convRepo ConversationRepository, carRepo CarRepository, leadRepo LeadRepository) *Bot {
	return &Bot{
		provider: provider,
		convRepo: convRepo,
		carRepo:  carRepo,
		leadRepo: leadRepo,
	}
}

// ProcessMessage processes incoming message and returns bot response
func (b *Bot) ProcessMessage(ctx context.Context, tenantID int, senderPhone, messageText string, isSales bool) (string, error) {
	slog.Info("processing message", "tenant_id", tenantID, "sender", senderPhone, "is_sales", isSales)

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

	case "createLead":
		phoneNumber, ok := arguments["phone_number"].(string)
		if !ok {
			return nil, fmt.Errorf("invalid phone_number")
		}

		name, _ := arguments["name"].(string)

		var carID *int
		if cid, ok := arguments["interested_car_id"].(float64); ok {
			id := int(cid)
			carID = &id
		}

		leadID, err := b.leadRepo.Create(ctx, phoneNumber, name, carID)
		if err != nil {
			return nil, fmt.Errorf("failed to create lead: %w", err)
		}

		return map[string]interface{}{
			"lead_id": leadID,
			"message": "Lead berhasil dibuat",
		}, nil

	default:
		return nil, fmt.Errorf("unknown function: %s", functionName)
	}
}

// buildSystemPrompt creates system prompt based on user type
func (b *Bot) buildSystemPrompt(isSales bool) string {
	if isSales {
		return `Anda adalah asisten penjualan mobil untuk tim sales internal.

Anda membantu sales mengakses informasi inventory mobil dengan cepat.

Kemampuan Anda:
- Mencari mobil berdasarkan kriteria (merek, harga, transmisi, dll)
- Memberikan detail lengkap mobil termasuk spesifikasi
- Membuat catatan lead dari customer

Gunakan bahasa Indonesia yang profesional dan efisien.`
	}

	return `Anda adalah asisten penjualan mobil yang ramah dan membantu customer mencari mobil yang sesuai.

Kemampuan Anda:
- Mencari mobil berdasarkan budget dan preferensi customer
- Memberikan informasi detail mobil termasuk harga dan spesifikasi
- Membantu customer membandingkan pilihan
- Mengatur jadwal untuk melihat mobil atau test drive

Gunakan bahasa Indonesia yang ramah, natural, dan helpful.
Pahami istilah automotive Indonesia seperti: matic (automatic), bensin (gasoline), OTR (On The Road price).`
}

// GetAvailableFunctions returns functions available to LLM
func (b *Bot) GetAvailableFunctions(isSales bool) []Function {
	functions := []Function{
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
	}

	if isSales {
		// Sales get additional function to create leads
		functions = append(functions, Function{
			Name:        "createLead",
			Description: "Buat lead baru dari customer yang tertarik",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"phone_number": map[string]interface{}{
						"type":        "string",
						"description": "Nomor telepon customer",
					},
					"name": map[string]interface{}{
						"type":        "string",
						"description": "Nama customer",
					},
					"interested_car_id": map[string]interface{}{
						"type":        "integer",
						"description": "ID mobil yang diminati",
					},
				},
				"required": []string{"phone_number"},
			},
		})
	}

	return functions
}
