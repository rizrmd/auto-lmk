package llm

import (
	"context"

	"github.com/riz/auto-lmk/internal/model"
	"github.com/riz/auto-lmk/internal/repository"
)

// ConversationRepoAdapter adapts repository.ConversationRepository to bot's interface
type ConversationRepoAdapter struct {
	repo *repository.ConversationRepository
}

func NewConversationRepoAdapter(repo *repository.ConversationRepository) *ConversationRepoAdapter {
	return &ConversationRepoAdapter{repo: repo}
}

func (a *ConversationRepoAdapter) GetOrCreate(ctx context.Context, senderPhone string, isSales bool) (*Conversation, error) {
	conv, err := a.repo.GetOrCreate(ctx, senderPhone, isSales)
	if err != nil {
		return nil, err
	}
	return &Conversation{ID: conv.ID}, nil
}

func (a *ConversationRepoAdapter) AddMessage(ctx context.Context, conversationID int, senderPhone, messageText, direction string) error {
	return a.repo.AddMessage(ctx, conversationID, senderPhone, messageText, direction)
}

func (a *ConversationRepoAdapter) GetMessages(ctx context.Context, conversationID int, limit int) ([]*BotMessage, error) {
	msgs, err := a.repo.GetMessages(ctx, conversationID, limit)
	if err != nil {
		return nil, err
	}

	result := make([]*BotMessage, len(msgs))
	for i, msg := range msgs {
		result[i] = &BotMessage{
			SenderPhone: msg.SenderPhone,
			MessageText: msg.MessageText,
			Direction:   msg.Direction,
		}
	}
	return result, nil
}

// CarRepoAdapter adapts service.CarService to bot's interface
type CarRepoAdapter struct {
	car *repository.CarRepository
}

func NewCarRepoAdapter(carRepo *repository.CarRepository) *CarRepoAdapter {
	return &CarRepoAdapter{car: carRepo}
}

func (a *CarRepoAdapter) SearchCarsForBot(ctx context.Context, filters map[string]interface{}) (interface{}, error) {
	// Smart filter conversions
	if trans, ok := filters["transmission"].(string); ok {
		if trans == "matic" {
			filters["transmission"] = "automatic"
		}
	}

	return a.car.List(ctx, filters)
}

func (a *CarRepoAdapter) GetCarWithDetails(ctx context.Context, carID int) (interface{}, error) {
	return a.car.GetByID(ctx, carID)
}

// LeadRepoAdapter adapts repository.LeadRepository to bot's interface
type LeadRepoAdapter struct {
	repo *repository.LeadRepository
}

func NewLeadRepoAdapter(repo *repository.LeadRepository) *LeadRepoAdapter {
	return &LeadRepoAdapter{repo: repo}
}

func (a *LeadRepoAdapter) Create(ctx context.Context, phoneNumber, name string, interestedCarID *int) (int, error) {
	var namePtr *string
	if name != "" {
		namePtr = &name
	}

	req := &model.CreateLeadRequest{
		PhoneNumber:     phoneNumber,
		Name:            namePtr,
		InterestedCarID: interestedCarID,
	}

	lead, err := a.repo.Create(ctx, req)
	if err != nil {
		return 0, err
	}

	return lead.ID, nil
}
