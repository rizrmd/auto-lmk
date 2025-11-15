package whatsapp

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"sync"

	_ "github.com/lib/pq"
	"github.com/riz/auto-lmk/internal/repository"
	"github.com/skip2/go-qrcode"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
)

// Client manages WhatsApp connections for multiple tenants
type Client struct {
	salesRepo     *repository.SalesRepository
	clients       map[int]*whatsmeow.Client // tenant_id -> whatsmeow client
	container     *sqlstore.Container
	mu            sync.RWMutex
	messageHandler MessageHandler
}

// MessageHandler processes incoming messages
type MessageHandler func(ctx context.Context, tenantID int, senderPhone, messageText string) error

// NewClient creates a new WhatsApp client manager
func NewClient(salesRepo *repository.SalesRepository, dbURL string) (*Client, error) {
	// Create store for WhatsApp session data
	dbLog := waLog.Stdout("Database", "INFO", true)

	// Connect to database for session storage
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	container := sqlstore.NewWithDB(db, "postgres", dbLog)
	if err := container.Upgrade(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to upgrade sqlstore: %w", err)
	}

	return &Client{
		salesRepo: salesRepo,
		clients:   make(map[int]*whatsmeow.Client),
		container: container,
	}, nil
}

// SetMessageHandler sets the handler for incoming messages
func (c *Client) SetMessageHandler(handler MessageHandler) {
	c.messageHandler = handler
}

// PairTenant initiates WhatsApp pairing for a tenant
func (c *Client) PairTenant(ctx context.Context, tenantID int, phoneNumber string) (string, error) {
	slog.Info("initiating WhatsApp pairing", "tenant_id", tenantID, "phone", phoneNumber)

	c.mu.Lock()
	defer c.mu.Unlock()

	// Check if already paired
	if client, exists := c.clients[tenantID]; exists && client.IsConnected() {
		return "", fmt.Errorf("tenant already paired and connected")
	}

	// Create device store
	deviceStore := c.container.NewDevice()
	clientLog := waLog.Stdout("Client", "INFO", true)

	// Create WhatsApp client
	waClient := whatsmeow.NewClient(deviceStore, clientLog)

	// Generate QR code for pairing
	qrChan, err := waClient.GetQRChannel(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get QR channel: %w", err)
	}

	err = waClient.Connect()
	if err != nil {
		return "", fmt.Errorf("failed to connect: %w", err)
	}

	// Wait for QR code
	var qrCode string
	for evt := range qrChan {
		if evt.Event == "code" {
			// Generate QR code image
			qrCodePath := fmt.Sprintf("/tmp/qr_%d.png", tenantID)
			err := qrcode.WriteFile(evt.Code, qrcode.Medium, 256, qrCodePath)
			if err != nil {
				slog.Error("failed to generate QR code image", "error", err)
			}

			qrCode = evt.Code
			slog.Info("QR code generated", "tenant_id", tenantID, "path", qrCodePath)
			break
		} else if evt.Event == "success" {
			slog.Info("QR code scan successful", "tenant_id", tenantID)
			// Store client
			c.clients[tenantID] = waClient
			// Set up message handler
			waClient.AddEventHandler(c.createEventHandler(tenantID))
			return "PAIRED_SUCCESSFULLY", nil
		}
	}

	return qrCode, nil
}

// SendMessage sends a WhatsApp message
func (c *Client) SendMessage(tenantID int, recipientPhone, message string) error {
	slog.Info("sending WhatsApp message", "tenant_id", tenantID, "recipient", recipientPhone)

	c.mu.RLock()
	client, exists := c.clients[tenantID]
	c.mu.RUnlock()

	if !exists || !client.IsConnected() {
		return fmt.Errorf("tenant not connected to WhatsApp")
	}

	// Parse phone number to JID
	jid, err := parsePhoneNumber(recipientPhone)
	if err != nil {
		return fmt.Errorf("invalid phone number: %w", err)
	}

	// Send message
	_, err = client.SendMessage(context.Background(), jid, &waE2E.Message{
		Conversation: proto.String(message),
	})

	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	slog.Info("message sent successfully", "tenant_id", tenantID, "recipient", recipientPhone)
	return nil
}

// SendImage sends a WhatsApp image
func (c *Client) SendImage(tenantID int, recipientPhone, imagePath, caption string) error {
	slog.Info("sending WhatsApp image", "tenant_id", tenantID, "recipient", recipientPhone, "path", imagePath)

	c.mu.RLock()
	client, exists := c.clients[tenantID]
	c.mu.RUnlock()

	if !exists || !client.IsConnected() {
		return fmt.Errorf("tenant not connected to WhatsApp")
	}

	// Parse phone number to JID
	jid, err := parsePhoneNumber(recipientPhone)
	if err != nil {
		return fmt.Errorf("invalid phone number: %w", err)
	}

	// Read image file
	imageData, err := os.ReadFile(imagePath)
	if err != nil {
		return fmt.Errorf("failed to read image: %w", err)
	}

	// Upload image
	uploaded, err := client.Upload(context.Background(), imageData, whatsmeow.MediaImage)
	if err != nil {
		return fmt.Errorf("failed to upload image: %w", err)
	}

	// Send image message
	imageMsg := &waE2E.ImageMessage{
		Caption:        proto.String(caption),
		URL:            proto.String(uploaded.URL),
		DirectPath:     proto.String(uploaded.DirectPath),
		MediaKey:       uploaded.MediaKey,
		Mimetype:       proto.String("image/jpeg"),
		FileEncSHA256:  uploaded.FileEncSHA256,
		FileSHA256:     uploaded.FileSHA256,
		FileLength:     proto.Uint64(uint64(len(imageData))),
	}

	_, err = client.SendMessage(context.Background(), jid, &waE2E.Message{
		ImageMessage: imageMsg,
	})

	if err != nil {
		return fmt.Errorf("failed to send image: %w", err)
	}

	slog.Info("image sent successfully", "tenant_id", tenantID, "recipient", recipientPhone)
	return nil
}

// IsConnected checks if tenant's WhatsApp is connected
func (c *Client) IsConnected(tenantID int) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	client, exists := c.clients[tenantID]
	if !exists {
		return false
	}

	return client.IsConnected()
}

// Disconnect disconnects a tenant's WhatsApp
func (c *Client) Disconnect(tenantID int) error {
	slog.Info("disconnecting WhatsApp", "tenant_id", tenantID)

	c.mu.Lock()
	defer c.mu.Unlock()

	client, exists := c.clients[tenantID]
	if !exists {
		return fmt.Errorf("tenant not found")
	}

	// Disconnect
	client.Disconnect()

	// Remove from map
	delete(c.clients, tenantID)

	slog.Info("WhatsApp disconnected", "tenant_id", tenantID)
	return nil
}

// createEventHandler creates an event handler for incoming messages
func (c *Client) createEventHandler(tenantID int) func(interface{}) {
	return func(evt interface{}) {
		switch v := evt.(type) {
		case *events.Message:
			// Handle incoming message
			if !v.Info.IsFromMe && v.Message.GetConversation() != "" {
				senderPhone := v.Info.Sender.User // Phone number without country code
				messageText := v.Message.GetConversation()

				slog.Info("received message", "tenant_id", tenantID, "sender", senderPhone, "text", messageText)

				// Call message handler if set
				if c.messageHandler != nil {
					ctx := context.Background()
					if err := c.messageHandler(ctx, tenantID, senderPhone, messageText); err != nil {
						slog.Error("failed to handle message", "error", err)
					}
				}
			}
		case *events.Connected:
			slog.Info("WhatsApp connected", "tenant_id", tenantID)
		case *events.Disconnected:
			slog.Warn("WhatsApp disconnected", "tenant_id", tenantID)
		}
	}
}

// parsePhoneNumber converts phone number to WhatsApp JID
func parsePhoneNumber(phone string) (types.JID, error) {
	// Remove any non-numeric characters
	cleaned := ""
	for _, r := range phone {
		if r >= '0' && r <= '9' {
			cleaned += string(r)
		}
	}

	if len(cleaned) == 0 {
		return types.JID{}, fmt.Errorf("invalid phone number")
	}

	// If number doesn't start with country code, assume Indonesia (+62)
	if !strings.HasPrefix(cleaned, "62") && len(cleaned) < 12 {
		// Remove leading 0 if present
		if strings.HasPrefix(cleaned, "0") {
			cleaned = cleaned[1:]
		}
		cleaned = "62" + cleaned
	}

	return types.NewJID(cleaned, types.DefaultUserServer), nil
}
