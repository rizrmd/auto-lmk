# WhatsApp Bot Testing Guide

**Date:** 2025-11-15
**Status:** Ready for Testing
**Provider:** Z.AI (GLM-4-Flash)

## Overview

The Auto LMK WhatsApp bot is now fully integrated and ready for testing. This guide will help you set up and test the bot with Z.AI's LLM provider.

## Prerequisites

- [x] PostgreSQL database running
- [x] Z.AI API key configured
- [x] WhatsApp Business Account or personal WhatsApp (for testing)
- [x] Server/local machine running the Auto LMK API

## Architecture

```
WhatsApp Message
    ↓
Whatsmeow Client
    ↓
Message Handler
    ↓
WhatsApp Service
    ↓
LLM Bot (with Z.AI)
    ↓
Function Calls (searchCars, getCarDetails, createLead)
    ↓
Response back to WhatsApp
```

## Setup Instructions

### 1. Environment Configuration

Ensure your `.env` file has the following:

```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=autolmk
DB_PASSWORD=autolmk_dev_password
DB_NAME=autolmk
DB_SSLMODE=disable

# LLM Configuration
LLM_PROVIDER=zai
LLM_API_KEY=93ac6b4e9c1c49b4b64fed617669e569.5nfnaoMbbNaKZ26I
LLM_MODEL=glm-4-flash
ZAI_ENDPOINT=https://api.z.ai/api/coding/paas/v4

# WhatsApp
WHATSAPP_SESSION_PATH=./whatsapp_sessions
```

### 2. Database Setup

Run migrations to ensure all tables exist:

```bash
# From project root
psql -h localhost -U autolmk -d autolmk -f migrations/001_init.sql
```

### 3. Create Test Data

#### Create a Tenant

```bash
curl -X POST http://localhost:8080/api/admin/tenants \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Showroom",
    "domain": "test.localhost",
    "whatsapp_number": "6281234567890",
    "contact_email": "test@example.com",
    "status": "active"
  }'
```

#### Add Test Cars

```bash
# Set tenant context
export TENANT_ID=1

curl -X POST http://localhost:8080/api/cars \
  -H "Content-Type: application/json" \
  -H "X-Tenant-Domain: test.localhost" \
  -d '{
    "brand": "Toyota",
    "model": "Avanza",
    "year": 2023,
    "price": 250000000,
    "mileage": 5000,
    "transmission": "automatic",
    "fuel_type": "bensin",
    "engine_cc": 1500,
    "seats": 7,
    "color": "Silver",
    "description": "Toyota Avanza 2023 matic, kondisi sangat baik",
    "status": "available"
  }'

curl -X POST http://localhost:8080/api/cars \
  -H "Content-Type: application/json" \
  -H "X-Tenant-Domain: test.localhost" \
  -d '{
    "brand": "Honda",
    "model": "CR-V",
    "year": 2022,
    "price": 450000000,
    "mileage": 10000,
    "transmission": "automatic",
    "fuel_type": "bensin",
    "engine_cc": 1500,
    "seats": 5,
    "color": "White",
    "description": "Honda CR-V Turbo Prestige, full option",
    "status": "available"
  }'
```

### 4. Start the Server

```bash
go run cmd/api/main.go
```

Expected output:
```
Auto LMK - Multi-Tenant Car Sales Platform
===========================================
INFO starting application env=development port=8080
INFO database connection successful max_open_conns=25 max_idle_conns=5
INFO LLM provider initialized provider=zai model=glm-4-flash
INFO WhatsApp client initialized
INFO WhatsApp bot service initialized
INFO server starting port=8080
```

### 5. Pair WhatsApp

*Note: WhatsApp pairing will be implemented via admin panel. For now, pairing happens programmatically.*

The QR code will be generated at `/tmp/qr_<tenant_id>.png`. Scan this with WhatsApp to pair.

## Test Scenarios

### Scenario 1: Customer Inquiry (Bahasa Indonesia)

**Test Case:** Customer asks about available cars

**WhatsApp Message:**
```
Ada mobil apa aja?
```

**Expected Bot Response:**
- Bot calls `searchCars` function with empty filters
- Returns list of available cars
- Response in natural Bahasa Indonesia

**Example Response:**
```
Hai! Kami punya beberapa mobil yang tersedia:

1. Toyota Avanza 2023 - Rp 250.000.000
   Matic, bensin, 7 seats, warna Silver

2. Honda CR-V 2022 - Rp 450.000.000
   Matic, bensin, 5 seats, warna White

Ada yang menarik perhatian Anda?
```

### Scenario 2: Filtered Search (Budget)

**WhatsApp Message:**
```
Toyota budget 300 juta
```

**Expected:**
- Bot understands: brand=Toyota, max_price=300000000
- Calls `searchCars` with filters
- Returns matching results

**Example Response:**
```
Saya menemukan Toyota dalam budget Anda:

1. Toyota Avanza 2023 - Rp 250.000.000
   Kondisi sangat baik, kilometer rendah (5.000 km)
   Matic, bensin, 7 kursi

Mau info lebih detail tentang Avanza ini?
```

### Scenario 3: Indonesian Automotive Terms

**WhatsApp Message:**
```
Cari mobil matic, bensin, 2020 ke atas
```

**Expected:**
- Bot understands "matic" = automatic transmission
- Filters: transmission=automatic, fuel_type=bensin, min_year=2020
- Returns matching cars

### Scenario 4: Car Details Request

**WhatsApp Message:**
```
Yang CR-V detail nya dong
```

**Expected:**
- Bot calls `getCarDetails` for Honda CR-V
- Returns full specifications
- Mentions photos will be sent (if implemented)

### Scenario 5: Lead Creation (Sales)

**Prerequisites:** Register a sales phone number first

```sql
INSERT INTO sales (tenant_id, phone_number, name, status)
VALUES (1, '6281234567890', 'Test Sales', 'active');
```

**WhatsApp Message from Sales:**
```
Ada customer mau beli Avanza, nama Budi 08123456789
```

**Expected:**
- Bot recognizes sender as sales user
- Calls `createLead` function
- Creates lead in database
- Confirms lead creation

**Example Response:**
```
Lead berhasil dibuat!

Customer: Budi
No HP: 08123456789
Minat: Toyota Avanza 2023
Status: New

Saya sudah simpan di sistem.
```

### Scenario 6: Multi-turn Conversation

**Turn 1:**
```
Ada mobil apa aja?
```

**Turn 2:**
```
Yang SUV aja
```

**Turn 3:**
```
Yang paling murah berapa?
```

**Expected:**
- Bot maintains conversation context
- Progressive filtering based on user input
- Understands "Yang" refers to previous results

## Testing Checklist

### Basic Functionality
- [ ] Server starts without errors
- [ ] WhatsApp client initializes
- [ ] QR code generated successfully
- [ ] WhatsApp pairs successfully
- [ ] Bot receives messages

### LLM Integration
- [ ] Z.AI API connection successful
- [ ] Indonesian language responses are natural
- [ ] Function calling works correctly
- [ ] Context maintained across messages

### Function Execution
- [ ] `searchCars` - Returns correct cars based on filters
- [ ] `getCarDetails` - Returns car details
- [ ] `createLead` - Creates lead successfully (sales only)

### Edge Cases
- [ ] No results found - graceful message
- [ ] Invalid car ID - error handling
- [ ] Unclear query - bot asks for clarification
- [ ] Large result set - pagination or summary

### Performance
- [ ] Response time < 5 seconds
- [ ] Concurrent messages handled correctly
- [ ] Database queries optimized
- [ ] No memory leaks after prolonged usage

## Monitoring

### Logs to Watch

```bash
# Follow application logs
tail -f /var/log/autolmk/app.log

# Watch for these log patterns:
# - "received message" - Incoming WhatsApp messages
# - "processing message" - LLM bot processing
# - "function call requested" - LLM calling a function
# - "function execution failed" - Function errors
# - "message sent successfully" - Response sent
```

### Database Monitoring

```sql
-- Check recent conversations
SELECT * FROM conversations ORDER BY updated_at DESC LIMIT 10;

-- Check recent messages
SELECT * FROM messages ORDER BY created_at DESC LIMIT 20;

-- Check leads created
SELECT * FROM leads ORDER BY created_at DESC LIMIT 10;

-- Check LLM usage (if logging implemented)
SELECT COUNT(*), AVG(response_time_ms) FROM llm_requests WHERE created_at > NOW() - INTERVAL '1 hour';
```

## Common Issues

### Issue: WhatsApp not receiving messages

**Symptoms:**
- Bot logs show "message sent successfully"
- WhatsApp doesn't receive the message

**Solutions:**
1. Check WhatsApp connection status
2. Verify phone number format (must include country code)
3. Check if WhatsApp session is still valid
4. Re-pair WhatsApp if needed

### Issue: LLM not responding

**Symptoms:**
- Messages received but no response
- Error logs show "LLM call failed"

**Solutions:**
1. Verify Z.AI API key is correct
2. Check Z.AI API status/limits
3. Review Z.AI endpoint configuration
4. Check network connectivity

### Issue: Function calls failing

**Symptoms:**
- Bot responds but with error messages
- Logs show "function execution failed"

**Solutions:**
1. Check database connection
2. Verify tenant context is set correctly
3. Review function argument parsing
4. Check data exists (cars, leads, etc.)

### Issue: Indonesian responses are poor quality

**Symptoms:**
- Responses in English or mixed
- Unnatural Indonesian phrasing
- Missing automotive terms

**Solutions:**
1. Adjust system prompts in `internal/llm/bot.go`
2. Add more Indonesian examples to prompts
3. Test different temperature settings
4. Consider switching to different model (GLM-4 instead of GLM-4-Flash)

## Performance Metrics

### Target Metrics

- **Response Time:** < 5 seconds (end-to-end)
- **LLM Latency:** < 2 seconds
- **Function Execution:** < 500ms
- **Database Queries:** < 100ms

### Cost Estimation

**Z.AI GLM-4-Flash Pricing** (estimated):
- Input: ~$0.50/1M tokens
- Output: ~$1.50/1M tokens

**Average Conversation:**
- ~200 tokens per message (input + output)
- ~7 messages per conversation
- = ~1,400 tokens per conversation

**Monthly Cost** (10 tenants, 15 conversations/day each):
- 10 × 15 × 7 × 1400 × 30 = 44,100,000 tokens (~44M tokens/month)
- Estimated cost: ~$44-66/month

## Next Steps

After successful testing:

1. **Production Deployment**
   - Set up production server
   - Configure production database
   - Set up SSL/HTTPS
   - Configure production WhatsApp numbers

2. **Monitoring Setup**
   - Set up error tracking (Sentry)
   - Configure uptime monitoring
   - Set up alerting for critical errors

3. **Optimization**
   - Cache common queries
   - Optimize database indexes
   - Implement rate limiting
   - Add response templates

4. **Feature Enhancements**
   - Photo sending implementation
   - Multi-language support
   - Advanced search filters
   - Analytics dashboard

## Support

For issues or questions:
- Check logs in `/var/log/autolmk/`
- Review configuration in `.env`
- Consult main documentation in `README.md`
- Review code in `internal/whatsapp/` and `internal/llm/`

## Success Criteria

Bot is production-ready when:
- ✅ All test scenarios pass
- ✅ No critical errors in logs
- ✅ Response time < 5 seconds consistently
- ✅ Z.AI costs within budget
- ✅ Indonesian responses are natural and correct
- ✅ Function calls execute successfully
- ✅ Multi-turn conversations work correctly
- ✅ Sales and customer modes both functional

---

**Last Updated:** 2025-11-15
**Version:** 1.0 (Initial WhatsApp Integration)
