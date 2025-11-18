package service

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/riz/auto-lmk/internal/llm"
	"github.com/riz/auto-lmk/internal/repository"
	"github.com/riz/auto-lmk/internal/whatsapp"
)

// WhatsAppService orchestrates WhatsApp bot functionality
type WhatsAppService struct {
	waClient   *whatsapp.Client
	bot        *llm.Bot
	salesRepo  *repository.SalesRepository
	convRepo   *repository.ConversationRepository
	carService *CarService
}

func NewWhatsAppService(
	waClient *whatsapp.Client,
	bot *llm.Bot,
	salesRepo *repository.SalesRepository,
	convRepo *repository.ConversationRepository,
	carService *CarService,
) *WhatsAppService {
	return &WhatsAppService{
		waClient:   waClient,
		bot:        bot,
		salesRepo:  salesRepo,
		convRepo:   convRepo,
		carService: carService,
	}
}

// ProcessIncomingMessage handles incoming WhatsApp message
func (s *WhatsAppService) ProcessIncomingMessage(ctx context.Context, tenantID int, senderPhone, messageText, messageType, mediaURL string) error {
	slog.Info("processing WhatsApp message", "tenant_id", tenantID, "sender", senderPhone, "type", messageType)

	// 1. Check if sender is sales
	isSales, err := s.salesRepo.IsSales(tenantID, senderPhone)
	if err != nil {
		slog.Error("failed to check sales status", "error", err)
		isSales = false // Default to customer
	}

	// 2. Handle image upload (Story 6.3)
	if messageType == "image" {
		if !isSales {
			// Customer tried to upload photo - reject politely (Story 6.7)
			response := "Terima kasih! Untuk upload mobil, silakan hubungi sales team kami. Saya dapat membantu Anda mencari mobil yang tersedia. Ada yang bisa saya bantu?"
			slog.Info("customer upload attempt blocked", "phone", senderPhone, "tenant_id", tenantID)

			if err := s.waClient.SendMessage(tenantID, senderPhone, response); err != nil {
				return fmt.Errorf("failed to send rejection message: %w", err)
			}
			return nil
		}

		// Sales uploaded photo - save it
		photoPath, err := s.saveCarPhoto(ctx, tenantID, mediaURL)
		if err != nil {
			slog.Error("failed to save photo", "error", err)
			errMsg := "Maaf, gagal menyimpan foto. Silakan coba lagi."
			s.waClient.SendMessage(tenantID, senderPhone, errMsg)
			return err
		}

		// Add to pending photos
		count := s.bot.AddPendingPhoto(senderPhone, photoPath)

		// Send confirmation
		var response string
		if count == 1 {
			response = "Foto diterima! Upload lebih banyak foto atau ketik detail mobil.\nContoh: Toyota Avanza 2020 185juta AT Bensin"
		} else if count < 5 {
			response = fmt.Sprintf("Foto %d diterima! Total %d foto.", count, count)
		} else {
			response = "Maksimum 5 foto tercapai. Silakan ketik detail mobil sekarang."
		}

		if err := s.waClient.SendMessage(tenantID, senderPhone, response); err != nil {
			return fmt.Errorf("failed to send confirmation: %w", err)
		}

		return nil
	}

	// 3. Get or create conversation (for text messages)
	conversation, err := s.convRepo.GetOrCreate(ctx, senderPhone, isSales)
	if err != nil {
		return fmt.Errorf("failed to get conversation: %w", err)
	}

	// 4. Store incoming message
	if err := s.convRepo.AddMessage(ctx, conversation.ID, senderPhone, messageText, "inbound"); err != nil {
		slog.Error("failed to store message", "error", err)
	}

	// 5. Process with LLM bot
	response, err := s.bot.ProcessMessage(ctx, tenantID, senderPhone, messageText, isSales)
	if err != nil {
		slog.Error("failed to process with bot", "error", err)
		response = "Maaf, saya sedang mengalami gangguan. Silakan coba lagi nanti."
	}

	// 6. Store bot response
	if err := s.convRepo.AddMessage(ctx, conversation.ID, "BOT", response, "outbound"); err != nil {
		slog.Error("failed to store bot response", "error", err)
	}

	// 7. Send via WhatsApp
	if err := s.waClient.SendMessage(tenantID, senderPhone, response); err != nil {
		return fmt.Errorf("failed to send WhatsApp message: %w", err)
	}

	// 8. Send pending images if any
	pendingImages := s.bot.GetPendingImages(senderPhone)
	if len(pendingImages) > 0 {
		slog.Info("sending car images", "tenant_id", tenantID, "sender", senderPhone, "count", len(pendingImages))

		for i, imagePath := range pendingImages {
			caption := ""
			if i == 0 {
				// Add caption to first image
				caption = "Foto mobil yang Anda minta"
			}

			err := s.waClient.SendImage(tenantID, senderPhone, imagePath, caption)
			if err != nil {
				slog.Error("failed to send image", "error", err, "path", imagePath)
				// Continue sending other images even if one fails
			}
		}

		// Clear pending images after sending
		s.bot.ClearPendingImages(senderPhone)
	}

	return nil
}

// saveCarPhoto saves an uploaded car photo to disk
func (s *WhatsAppService) saveCarPhoto(ctx context.Context, tenantID int, mediaURL string) (string, error) {
	// For simplified implementation, create a placeholder photo
	// In production, you would:
	// 1. Download from mediaURL using s.waClient.DownloadMedia()
	// 2. Validate format and size
	// 3. Generate thumbnail
	// 4. Save both original and thumbnail

	// Create upload directory if not exists
	uploadDir := filepath.Join("static", "uploads", "cars", fmt.Sprintf("%d", tenantID))
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %w", err)
	}

	// Generate filename with timestamp
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("%d_car.jpg", timestamp)
	photoPath := filepath.Join(uploadDir, filename)

	// For now, create a placeholder file
	// In production, write actual image data from WhatsApp
	if err := os.WriteFile(photoPath, []byte("placeholder image data"), 0644); err != nil {
		return "", fmt.Errorf("failed to save photo: %w", err)
	}

	slog.Info("photo saved", "path", photoPath, "tenant_id", tenantID)

	// Return relative path for database storage
	return fmt.Sprintf("/static/uploads/cars/%d/%s", tenantID, filename), nil
}

// ExecuteBotFunction executes function called by LLM
func (s *WhatsAppService) ExecuteBotFunction(ctx context.Context, functionName string, arguments map[string]interface{}) (interface{}, error) {
	slog.Info("executing bot function", "function", functionName, "args", arguments)

	switch functionName {
	case "searchCars":
		return s.carService.SearchCarsForBot(ctx, arguments)

	case "getCarDetails":
		carID, ok := arguments["car_id"].(float64)
		if !ok {
			return nil, fmt.Errorf("invalid car_id")
		}
		return s.carService.GetCarWithDetails(ctx, int(carID))

	default:
		return nil, fmt.Errorf("unknown function: %s", functionName)
	}
}
