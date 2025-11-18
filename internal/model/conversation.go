package model

import "time"

type Conversation struct {
	ID          int       `json:"id"`
	TenantID    int       `json:"tenant_id"`
	SenderPhone string    `json:"sender_phone"`
	IsSales     bool      `json:"is_sales"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Message struct {
	ID             int       `json:"id"`
	ConversationID int       `json:"conversation_id"`
	SenderPhone    string    `json:"sender_phone"`
	MessageText    string    `json:"message_text"`
	Direction      string    `json:"direction"` // inbound, outbound
	CreatedAt      time.Time `json:"created_at"`
}

// ConversationListItem represents conversation with last message info for list view
type ConversationListItem struct {
	ID            int       `json:"id"`
	PhoneNumber   string    `json:"phone_number"`
	IsSales       bool      `json:"is_sales"`
	LastMessage   string    `json:"last_message"`
	LastMessageAt time.Time `json:"last_message_at"`
	MessageCount  int       `json:"message_count"`
	CreatedAt     time.Time `json:"created_at"`
}

type Lead struct {
	ID               int       `json:"id"`
	TenantID         int       `json:"tenant_id"`
	PhoneNumber      string    `json:"phone_number"`
	Name             *string   `json:"name,omitempty"`
	InterestedCarID  *int      `json:"interested_car_id,omitempty"`
	ConversationID   *int      `json:"conversation_id,omitempty"`
	Status           string    `json:"status"` // new, contacted, converted, lost
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type CreateLeadRequest struct {
	PhoneNumber     string  `json:"phone_number"`
	Name            *string `json:"name,omitempty"`
	InterestedCarID *int    `json:"interested_car_id,omitempty"`
	ConversationID  *int    `json:"conversation_id,omitempty"`
}
