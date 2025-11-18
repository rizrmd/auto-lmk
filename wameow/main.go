package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	fmt.Println("Starting WhatsApp client...")
	
	// Database setup for session storage
	container, err := sqlstore.New("sqlite3", "file:session.db?_foreign_keys=on", nil)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Get device store
	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		log.Fatalf("Failed to get device store: %v", err)
	}

	// Create client
	clientLog := waLog.Stdout("Client", "INFO", true)
	client := whatsmeow.NewClient(deviceStore, clientLog)
	client.AddEventHandler(handler)

	// Check if we have an existing session
	if client.Store.ID != nil {
		fmt.Printf("Found existing session: %s\n", client.Store.ID.String())
		err = client.Connect()
		if err != nil {
			log.Fatalf("Failed to connect with existing session: %v", err)
		}
		fmt.Println("Connected successfully with existing session")
	} else {
		fmt.Println("No existing session found, generating QR code...")
		
		// Get QR channel BEFORE connecting
		qrChan, err := client.GetQRChannel(context.Background())
		if err != nil {
			log.Fatalf("Failed to get QR channel: %v", err)
		}
		
		fmt.Println("Connecting to WhatsApp...")
		err = client.Connect()
		if err != nil {
			log.Fatalf("Failed to connect: %v", err)
		}
		
		fmt.Println("Waiting for QR code...")
		
		// Wait for QR code events
		for evt := range qrChan {
			fmt.Printf("QR Event: %s\n", evt.Event)
			switch evt.Event {
			case "code":
				fmt.Println("QR Code generated:")
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
				fmt.Println("\n📱 Scan this QR code with your WhatsApp app")
				fmt.Println("⏱️ QR code will refresh every 60 seconds")
			case "success":
				fmt.Println("✅ Authentication successful!")
				goto authenticated
			case "timeout":
				fmt.Println("⏰ QR code timed out. Please try again.")
				return
			case "error":
				if evt.Error != nil {
					fmt.Printf("❌ QR code error: %v\n", evt.Error)
				}
				return
			case "err-client-outdated":
				fmt.Println("❌ Client version is outdated. Please update.")
				return
			case "err-scanned-without-multidevice":
				fmt.Println("❌ QR code was scanned but multidevice is not enabled on your phone.")
				fmt.Println("   Please enable multidevice in WhatsApp Settings > Linked Devices")
				return
			default:
				fmt.Printf("❓ Unknown QR event: %s\n", evt.Event)
			}
		}
	}

authenticated:
	fmt.Println("✅ Authenticated successfully!")
	fmt.Println("Client is running. Press Ctrl+C to stop.")

	// Wait for interrupt signal to gracefully shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	// Disconnect
	client.Disconnect()
	fmt.Println("Disconnected")
}

func handler(evt interface{}) {
	switch v := evt.(type) {
	case *whatsmeow.ConnectedEvent:
		fmt.Println("✅ Connected to WhatsApp!")
		if v.PairClientType == whatsmeow.PairClientWeb {
			fmt.Println("📱 Paired as web client")
		} else {
			fmt.Println("📱 Paired as mobile client")
		}

	case *whatsmeow.PairSuccessEvent:
		fmt.Println("✅ Pairing successful!")
		fmt.Printf("📱 Device ID: %s\n", v.ID.String())

	case *whatsmeow.PairErrorEvent:
		fmt.Printf("❌ Pairing error: %v\n", v.Error)

	case *whatsmeow.LoggedOutEvent:
		fmt.Println("👋 Logged out from WhatsApp")

	case *whatsmeow.AppStateSyncCompleteEvent:
		if len(v.Snapshot) > 0 {
			fmt.Printf("📊 App state sync complete with %d snapshot entries\n", len(v.Snapshot))
		} else {
			fmt.Println("📊 App state sync complete (no snapshot)")
		}

	case *whatsmeow.MessageEvent:
		fmt.Printf("💬 New message from %s: %s\n", v.Info.Sender.User, v.Message.GetConversation())

	case *whatsmeow.ReceiptEvent:
		fmt.Printf("✅ Receipt: %s from %s\n", v.Type.String(), v.MessageSender.String())

	default:
		fmt.Printf("📢 Event: %T\n", v)
	}
}