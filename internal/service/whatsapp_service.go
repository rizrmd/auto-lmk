package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/riz/auto-lmk/internal/llm"
	"github.com/riz/auto-lmk/internal/repository"
	"github.com/riz/auto-lmk/internal/whatsapp"
)

// WhatsAppService orchestrates WhatsApp bot functionality
type WhatsAppService struct {
	waClient     *whatsapp.Client
	bot          *llm.Bot
	salesRepo    *repository.SalesRepository
	convRepo     *repository.ConversationRepository
	carService   *CarService
	leadRepo     *repository.LeadRepository
}

func NewWhatsAppService(
	waClient *whatsapp.Client,
	bot *llm.Bot,
	salesRepo *repository.SalesRepository,
	convRepo *repository.ConversationRepository,
	carService *CarService,
	leadRepo *repository.LeadRepository,
) *WhatsAppService {
	return &WhatsAppService{
		waClient:   waClient,
		bot:        bot,
		salesRepo:  salesRepo,
		convRepo:   convRepo,
		carService: carService,
		leadRepo:   leadRepo,
	}
}

// ProcessIncomingMessage handles incoming WhatsApp message
func (s *WhatsAppService) ProcessIncomingMessage(ctx context.Context, tenantID int, senderPhone, messageText string) error {
	slog.Info("processing WhatsApp message", "tenant_id", tenantID, "sender", senderPhone)

	// 1. Check if sender is sales
	isSales, err := s.salesRepo.IsSales(tenantID, senderPhone)
	if err != nil {
		slog.Error("failed to check sales status", "error", err)
		isSales = false // Default to customer
	}

	// 2. Get or create conversation
	conversation, err := s.convRepo.GetOrCreate(ctx, senderPhone, isSales)
	if err != nil {
		return fmt.Errorf("failed to get conversation: %w", err)
	}

	// 3. Store incoming message
	if err := s.convRepo.AddMessage(ctx, conversation.ID, senderPhone, messageText, "inbound"); err != nil {
		slog.Error("failed to store message", "error", err)
	}

	// 4. Process with LLM bot
	response, err := s.bot.ProcessMessage(ctx, tenantID, senderPhone, messageText, isSales)
	if err != nil {
		slog.Error("failed to process with bot", "error", err)
		response = "Maaf, saya sedang mengalami gangguan. Silakan coba lagi nanti."
	}

	// 5. Store bot response
	if err := s.convRepo.AddMessage(ctx, conversation.ID, "BOT", response, "outbound"); err != nil {
		slog.Error("failed to store bot response", "error", err)
	}

	// 6. Send via WhatsApp
	if err := s.waClient.SendMessage(tenantID, senderPhone, response); err != nil {
		return fmt.Errorf("failed to send WhatsApp message: %w", err)
	}

	return nil
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

	case "createLead":
		// TODO: Implement lead creation
		return nil, fmt.Errorf("createLead not yet implemented")

	default:
		return nil, fmt.Errorf("unknown function: %s", functionName)
	}
}
