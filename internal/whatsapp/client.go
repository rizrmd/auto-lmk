package whatsapp

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"

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
	salesRepo      *repository.SalesRepository
	clients        map[int]*whatsmeow.Client // tenant_id -> whatsmeow client
	container      *sqlstore.Container
	mu             sync.RWMutex
	messageHandler MessageHandler
}

// MessageHandler processes incoming messages
type MessageHandler func(ctx context.Context, tenantID int, senderPhone, messageText, messageType, mediaURL string) error

// NewClient creates a new WhatsApp client manager
func NewClient(salesRepo *repository.SalesRepository, dbURL string) (*Client, error) {
	// Create store for WhatsApp session data
	dbLog := waLog.Stdout("Database", "INFO", true)

	// Connect to database for session storage
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Create custom session directory
	sessionDir := "/home/yopi/Projects/auto-lmk/whatsapp_sessions/.whatsapp"
	if err := os.MkdirAll(sessionDir, 0700); err != nil {
		slog.Warn("failed to create session directory", "path", sessionDir, "error", err)
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

// PairResult contains the result of pairing attempt
type PairResult struct {
	QRCode      string
	PhoneNumber string
	Status      string // "qr_generated" or "paired"
}

// PairTenant initiates WhatsApp pairing for a tenant
// NO PHONE NUMBER REQUIRED - WhatsApp auto-detects from device that scans QR
func (c *Client) PairTenant(ctx context.Context, tenantID int) (*PairResult, error) {
	slog.Info("[PAIRING-START] Initiating WhatsApp pairing", "tenant_id", tenantID, "step", "1-check-existing")

	c.mu.Lock()
	// Check if already paired
	if client, exists := c.clients[tenantID]; exists && client.IsConnected() {
		c.mu.Unlock()
		slog.Warn("[PAIRING-ABORT] Tenant already paired and connected", "tenant_id", tenantID)
		return nil, fmt.Errorf("tenant already paired and connected")
	}
	c.mu.Unlock()

	slog.Info("[PAIRING-STEP] Creating device store", "tenant_id", tenantID, "step", "2-create-store")
	// Create device store - whatsmeow akan menggunakan session directory utama
	deviceStore := c.container.NewDevice()
	clientLog := waLog.Stdout("Client", "INFO", true)

	slog.Info("[PAIRING-STEP] Creating WhatsApp client", "tenant_id", tenantID, "step", "3-create-client")
	// Create WhatsApp client
	waClient := whatsmeow.NewClient(deviceStore, clientLog)

	slog.Info("[PAIRING-STEP] Getting QR channel", "tenant_id", tenantID, "step", "4-get-qr-channel")
	// Generate QR code for pairing
	qrChan, err := waClient.GetQRChannel(ctx)
	if err != nil {
		slog.Error("[PAIRING-ERROR] Failed to get QR channel", "tenant_id", tenantID, "error", err)
		return nil, fmt.Errorf("failed to get QR channel: %w", err)
	}

	// Add additional logging for QR channel debugging
	slog.Info("[PAIRING-DEBUG] QR channel created", "tenant_id", tenantID, "channel_type", fmt.Sprintf("%T", qrChan))

	slog.Info("[PAIRING-DEBUG] Attempting to connect to WhatsApp server", "tenant_id", tenantID)
	err = waClient.Connect()
	if err != nil {
		slog.Error("[PAIRING-ERROR] Failed to connect to WhatsApp server", "tenant_id", tenantID, "error", err)
		return nil, fmt.Errorf("failed to connect: %w", err)
	}
	slog.Info("[PAIRING-SUCCESS] Connected to WhatsApp server", "tenant_id", tenantID, "step", "5-connected")

	// Add delay to ensure connection is stable
	time.Sleep(1 * time.Second)

	// Store the pending client temporarily
	slog.Info("[PAIRING-STEP] Storing client temporarily", "tenant_id", tenantID, "step", "6-store-client")
	c.mu.Lock()
	c.clients[tenantID] = waClient
	c.mu.Unlock()

	// Channel to communicate QR code from goroutine
	qrResult := make(chan *PairResult, 2) // Allow both QR and success results
	qrError := make(chan error, 1)

	slog.Info("[PAIRING-STEP] Starting QR event handler goroutine", "tenant_id", tenantID, "step", "7-start-goroutine")
	// Launch background goroutine to handle all QR events
	go func() {
		var qrCodeSent bool
		eventCount := 0

		slog.Info("[PAIRING-GOROUTINE] Started listening for QR events", "tenant_id", tenantID)

		for evt := range qrChan {
			eventCount++
			slog.Info("[PAIRING-EVENT] Received QR channel event", "tenant_id", tenantID, "event", evt.Event, "event_count", eventCount)

			switch evt.Event {
			case "code":
				slog.Info("[PAIRING-QR] Generating QR code", "tenant_id", tenantID, "qr_code_sent", qrCodeSent)
				// Generate QR code image
				qrCodePath := fmt.Sprintf("/tmp/qr_%d.png", tenantID)
				err := qrcode.WriteFile(evt.Code, qrcode.Medium, 256, qrCodePath)
				if err != nil {
					slog.Error("[PAIRING-ERROR] Failed to generate QR code image", "tenant_id", tenantID, "error", err)
				} else {
					slog.Info("[PAIRING-QR] QR code file created", "tenant_id", tenantID, "path", qrCodePath)
				}

				// Send QR code back to main thread (only first time)
				if !qrCodeSent {
					slog.Info("[PAIRING-QR] Sending QR code to main thread", "tenant_id", tenantID)
					qrResult <- &PairResult{
						QRCode: evt.Code,
						Status: "qr_generated",
					}
					qrCodeSent = true
					slog.Info("[PAIRING-QR] QR code sent successfully, waiting for scan", "tenant_id", tenantID)
				} else {
					slog.Info("[PAIRING-QR] Received additional QR code (refreshed)", "tenant_id", tenantID, "event_count", eventCount)
				}

			case "success":
				slog.Info("[PAIRING-SUCCESS] QR code scan successful!", "tenant_id", tenantID)

				// Get phone number from connected device
				phoneNumber := waClient.Store.ID.User
				slog.Info("[PAIRING-SUCCESS] WhatsApp paired successfully", "tenant_id", tenantID, "phone", phoneNumber)

				// Set up message handler
				slog.Info("[PAIRING-STEP] Setting up message event handler", "tenant_id", tenantID)
				waClient.AddEventHandler(c.createEventHandler(tenantID))

				slog.Info("[PAIRING-COMPLETE] Pairing completed successfully", "tenant_id", tenantID, "total_events", eventCount)

				// Send success result back to main thread
				qrResult <- &PairResult{
					PhoneNumber: phoneNumber,
					Status:      "paired",
				}
				return

			case "timeout":
				slog.Warn("[PAIRING-TIMEOUT] QR code expired", "tenant_id", tenantID, "event_count", eventCount)

			default:
				slog.Warn("[PAIRING-EVENT] Unknown QR event type", "tenant_id", tenantID, "event", evt.Event)
			}
		}

		// QR channel closed without success - this is normal behavior
		slog.Info("[PAIRING-INFO] QR channel closed (normal behavior)", "tenant_id", tenantID, "total_events", eventCount, "qr_code_sent", qrCodeSent)
		c.mu.Lock()
		delete(c.clients, tenantID)
		c.mu.Unlock()

		if !qrCodeSent {
			slog.Error("[PAIRING-ERROR] QR channel closed before any QR code was generated", "tenant_id", tenantID)
			qrError <- fmt.Errorf("QR channel closed unexpectedly")
		} else {
			slog.Info("[PAIRING-INFO] QR code sent, waiting for scan via Connected event", "tenant_id", tenantID)
		}
	}()

	// Wait for QR code generation (short timeout)
	slog.Info("[PAIRING-WAIT] Waiting for QR code generation", "tenant_id", tenantID, "timeout", "10s")
	qrCtx, qrCancel := context.WithTimeout(ctx, 10*time.Second)
	defer qrCancel()

	select {
	case result := <-qrResult:
		slog.Info("[PAIRING-RESULT] QR code received from goroutine", "tenant_id", tenantID, "status", result.Status)
		return result, nil
	case err := <-qrError:
		slog.Error("[PAIRING-ERROR] Error received from goroutine", "tenant_id", tenantID, "error", err)
		return nil, err
	case <-qrCtx.Done():
		slog.Error("[PAIRING-TIMEOUT] QR generation timeout", "tenant_id", tenantID, "error", qrCtx.Err())
		return nil, fmt.Errorf("QR generation timeout: %w", qrCtx.Err())
	}
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
		Caption:       proto.String(caption),
		URL:           proto.String(uploaded.URL),
		DirectPath:    proto.String(uploaded.DirectPath),
		MediaKey:      uploaded.MediaKey,
		Mimetype:      proto.String("image/jpeg"),
		FileEncSHA256: uploaded.FileEncSHA256,
		FileSHA256:    uploaded.FileSHA256,
		FileLength:    proto.Uint64(uint64(len(imageData))),
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

// DownloadMedia downloads media from WhatsApp (simplified placeholder for now)
func (c *Client) DownloadMedia(tenantID int, mediaURL string) ([]byte, error) {
	slog.Info("downloading media", "tenant_id", tenantID, "url", mediaURL)

	c.mu.RLock()
	client, exists := c.clients[tenantID]
	c.mu.RUnlock()

	if !exists || !client.IsConnected() {
		return nil, fmt.Errorf("tenant not connected to WhatsApp")
	}

	// For now, return placeholder
	// In production, you would:
	// 1. Download from WhatsApp servers using client.Download()
	// 2. Decrypt the media using MediaKey
	// 3. Return the decrypted data

	// Placeholder: Return empty data
	// Real implementation would use client.Download(imageMsg) with proper decryption
	return nil, fmt.Errorf("media download not yet implemented in simplified version")
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

// GetPhoneNumber returns the phone number of connected WhatsApp
func (c *Client) GetPhoneNumber(tenantID int) string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	client, exists := c.clients[tenantID]
	if !exists || !client.IsConnected() {
		return ""
	}

	if client.Store != nil && client.Store.ID != nil {
		return client.Store.ID.User
	}

	return ""
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
			if !v.Info.IsFromMe {
				senderPhone := v.Info.Sender.User // Phone number without country code

				// Handle text message
				if v.Message.GetConversation() != "" {
					messageText := v.Message.GetConversation()
					slog.Info("received text message", "tenant_id", tenantID, "sender", senderPhone, "text", messageText)

					// Call message handler if set
					if c.messageHandler != nil {
						ctx := context.Background()
						if err := c.messageHandler(ctx, tenantID, senderPhone, messageText, "text", ""); err != nil {
							slog.Error("failed to handle message", "error", err)
						}
					}
				}

				// Handle image message
				if v.Message.GetImageMessage() != nil {
					imageMsg := v.Message.GetImageMessage()
					slog.Info("received image message", "tenant_ID", tenantID, "sender", senderPhone)

					// Get image URL from WhatsApp
					mediaURL := imageMsg.GetURL()
					if mediaURL == "" && imageMsg.GetDirectPath() != "" {
						// Construct media URL if not directly provided
						mediaURL = imageMsg.GetDirectPath()
					}

					// Call message handler if set
					if c.messageHandler != nil {
						ctx := context.Background()
						caption := imageMsg.GetCaption()
						if err := c.messageHandler(ctx, tenantID, senderPhone, caption, "image", mediaURL); err != nil {
							slog.Error("failed to handle image", "error", err)
						}
					}
				}
			}
		case *events.Connected:
			slog.Info("[PAIRING-EVENT] WhatsApp connected event received", "tenant_id", tenantID)
			// This is where we detect successful pairing!
			phoneNumber := c.clients[tenantID].Store.ID.User
			slog.Info("[PAIRING-SUCCESS] Pairing success detected via Connected event", "tenant_id", tenantID, "phone", phoneNumber)

			// Update database with successful pairing
			// TODO: Add database update here

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
