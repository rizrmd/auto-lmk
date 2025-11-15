# WhatsApp Admin Panel - User Guide

**Last Updated:** 2025-11-15
**Version:** 1.0

## Overview

The WhatsApp Admin Panel allows tenant administrators to manage their WhatsApp bot connection directly from the web interface. No technical knowledge required!

## Accessing the Panel

1. Navigate to your admin dashboard
2. Click on "WhatsApp Settings" in the sidebar
3. Or go directly to: `https://your-domain.com/admin/whatsapp`

## Features

### 1. Connection Status

The panel displays real-time connection status:
- ✅ **Connected** - Your WhatsApp bot is active and receiving messages
- ❌ **Not Connected** - You need to pair your WhatsApp

### 2. Pairing WhatsApp

**Steps to connect:**

1. Click the "📱 Connect WhatsApp" button
2. A QR code will appear on screen
3. Open WhatsApp on your phone
4. Tap Menu (⋮) → "Linked Devices"
5. Tap "Link a Device"
6. Scan the QR code displayed on your screen
7. Wait for confirmation

**Important Notes:**
- QR code expires after 60 seconds - refresh if needed
- Keep the page open while scanning
- You can only link one WhatsApp number per tenant
- Use WhatsApp Business for best results

### 3. Disconnect WhatsApp

If you need to disconnect:

1. Click the "🔌 Disconnect" button
2. Confirm the action
3. Your WhatsApp will be unlinked

**When to disconnect:**
- Changing WhatsApp numbers
- Troubleshooting connection issues
- Security reasons (lost device)

### 4. Test Bot

Once connected, you can test your bot:

1. Enter a phone number (yours or a test number)
2. Click "📤 Send Test Message"
3. Check WhatsApp for the test message

**Test Message:**
```
Halo! Ini adalah pesan test dari Auto LMK Bot.
Bot WhatsApp Anda sudah berhasil terhubung! 🎉
```

## How the Bot Works

### For Customers

Customers can message your WhatsApp number and:
- Ask about available cars in natural Bahasa Indonesia
- Search by brand, price, features
- Get car details and photos
- Request test drives

**Example Customer Queries:**
```
"Ada Toyota budget 300 juta?"
"Cari mobil matic, bensin"
"Yang CR-V detail nya dong"
```

### For Sales Team

Registered sales team members get additional features:
- Access full inventory
- Create leads directly from chat
- View all car details

**Example Sales Commands:**
```
"List semua mobil di bawah 400 juta"
"Ada customer mau beli Avanza, nama Budi 0812345"
```

## Troubleshooting

### QR Code Won't Load

**Problem:** QR code displays "Loading..." indefinitely

**Solutions:**
1. Refresh the page
2. Check your internet connection
3. Try a different browser
4. Check server logs for errors

### WhatsApp Won't Connect

**Problem:** QR code scanned but status remains "Not Connected"

**Solutions:**
1. Wait 10-20 seconds (it takes time to establish connection)
2. Refresh the status by clicking 🔄
3. Try disconnecting and re-pairing
4. Ensure WhatsApp app is updated
5. Check if you have too many linked devices (max 4)

### Bot Not Responding

**Problem:** WhatsApp connected but bot doesn't respond to messages

**Solutions:**
1. Check connection status (should show ✅ Connected)
2. Send a test message to verify
3. Check LLM provider configuration in `.env`
4. Review server logs for errors
5. Verify tenant has car inventory

### Test Message Fails

**Problem:** "Failed to send message" error

**Solutions:**
1. Verify phone number format (include country code: 628xxx)
2. Check WhatsApp connection status
3. Ensure number exists and has WhatsApp
4. Check server logs

## API Endpoints

For advanced users or integrations:

### Get Status
```http
GET /api/admin/whatsapp/status
Headers: X-Tenant-Domain: your-domain.com
```

Response:
```json
{
  "tenant_id": 1,
  "is_connected": true,
  "phone_number": "6281234567890",
  "pairing_status": "paired"
}
```

### Initiate Pairing
```http
POST /api/admin/whatsapp/pair
Headers: X-Tenant-Domain: your-domain.com
```

Response:
```json
{
  "status": "pairing_initiated",
  "qr_code": "2@abc123...",
  "message": "Scan the QR code with WhatsApp to pair"
}
```

### Disconnect
```http
POST /api/admin/whatsapp/disconnect
Headers: X-Tenant-Domain: your-domain.com
```

### Send Test Message
```http
POST /api/admin/whatsapp/test
Headers: X-Tenant-Domain: your-domain.com
Content-Type: application/json

{
  "phone_number": "6281234567890",
  "message": "Test message (optional)"
}
```

## Best Practices

### Security
- Only share WhatsApp QR codes with authorized personnel
- Disconnect WhatsApp if device is lost or stolen
- Regularly review linked devices in WhatsApp settings
- Use WhatsApp Business for better security features

### Performance
- Keep connection active (avoid frequent disconnect/reconnect)
- Monitor bot response times
- Test bot regularly with sample queries
- Keep car inventory updated

### Maintenance
- Check connection status daily
- Test bot weekly
- Update WhatsApp app regularly
- Monitor message logs for issues

## Support

### Common Questions

**Q: Can I connect multiple WhatsApp numbers?**
A: One WhatsApp number per tenant. Contact support for multi-number setup.

**Q: What happens if connection drops?**
A: The system will attempt automatic reconnection. If it fails, re-pair via admin panel.

**Q: Can customers see I'm using a bot?**
A: The bot responds naturally in Bahasa Indonesia. You can customize messages in settings.

**Q: How much does it cost?**
A: WhatsApp integration is free. LLM usage costs ~$44-66/month for average usage.

### Getting Help

1. Check server logs: `/var/log/autolmk/app.log`
2. Review documentation: `/docs/whatsapp-bot-testing.md`
3. Contact technical support
4. Check GitHub issues: https://github.com/riz/auto-lmk/issues

## Updates & Changes

**Version 1.0 (2025-11-15)**
- Initial release
- QR code pairing
- Connection status monitoring
- Test message functionality

---

**Made with ❤️ by the Auto LMK Team**
