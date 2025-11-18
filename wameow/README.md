# WhatsApp Meow Go Examples

Complete Go examples for the `whatsmeow` library with QR code pairing and messaging capabilities.

## Files

- `main.go` - **Fixed main example with proper QR pairing**
- `simple_example.go` - Simplified version with cleaner QR handling
- `messaging_example.go` - Advanced messaging functions (text, media, polls, etc.)
- `interactive_example.go` - Interactive menu for testing features

## Quick Start

### 1. Main Example (Fixed)
```bash
go run main.go
```

Features:
- ✅ Proper QR code generation and display
- ✅ Session persistence
- ✅ Event handling with emojis
- ✅ Graceful shutdown
- ✅ Error handling for common issues

### 2. Simple Example
```bash
go run simple_example.go
```

Clean, minimal version focusing on core functionality.

### 3. Interactive Demo
```bash
go run interactive_example.go
```

Menu-driven interface for testing different features.

## Key Features Implemented

### Authentication
- QR code generation and terminal display
- Automatic QR refresh
- Session persistence with SQLite
- Error handling for common pairing issues
- Multi-device support detection

### Messaging Functions
- Send text messages
- Send replies with context
- Send messages with mentions
- Upload and send images
- Create polls
- Get contact/group info

### Event Handling
- Connection events
- Message events
- Receipt events
- Pairing success/error events
- App state sync events

## Dependencies

All examples use:
- `go.mau.fi/whatsmeow` - Main WhatsApp library
- `github.com/mdp/qrterminal/v3` - QR code display
- `github.com/mattn/go-sqlite3` - SQLite for session storage

## Troubleshooting

### QR Code Issues
- Ensure multidevice is enabled in WhatsApp Settings
- QR codes refresh every 60 seconds
- Check internet connection
- Try rescanning if timeout occurs

### Connection Issues
- Delete `session.db` to start fresh pairing
- Check firewall/proxy settings
- Ensure WhatsApp service is available

### Common Errors
- `err-client-outdated`: Update library version
- `err-scanned-without-multidevice`: Enable multidevice in WhatsApp

## Usage Examples

### Send Message
```go
targetJID := types.NewJID("1234567890", types.DefaultUserServer)
msg := &types.Message{
    Conversation: "Hello World!",
}
resp, err := client.SendMessage(context.Background(), targetJID, msg, nil)
```

### Upload Image
```go
uploadResp, err := client.Upload(context.Background(), imageData, whatsmeow.MediaImage)
msg := &types.Message{
    ImageMessage: &types.ImageMessage{
        Caption:   proto.String("Check this out!"),
        URL:       &uploadResp.URL,
        DirectPath: &uploadResp.DirectPath,
        // ... other fields
    },
}
client.SendMessage(context.Background(), targetJID, msg, nil)
```

### Create Poll
```go
pollMsg := client.BuildPollCreation(
    "Favorite language?",
    []string{"Go", "Python", "JavaScript"},
    1, // Max 1 selection
)
client.SendMessage(context.Background(), targetJID, pollMsg, nil)
```

## Security Notes

- Never share `session.db` files
- Store securely with appropriate permissions
- Use environment variables for sensitive data
- Follow WhatsApp's Terms of Service

## Next Steps

1. Run `main.go` to test basic connection
2. Try `simple_example.go` for clean implementation
3. Use `messaging_example.go` functions in your code
4. Customize event handlers for your use case