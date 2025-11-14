# API Testing Guide - Auto LMK Platform

**Created:** 2025-11-14
**Status:** Ready for Testing

## Prerequisites

1. **Start PostgreSQL:**
   ```bash
   docker-compose up -d
   ```

2. **Run Migrations:**
   ```bash
   make migrate-up
   # OR
   migrate -path migrations -database "postgresql://autolmk:autolmk_dev_password@localhost:5432/autolmk?sslmode=disable" up
   ```

3. **Seed Sample Data (Optional):**
   ```bash
   psql -U autolmk -d autolmk -h localhost < scripts/seed.sql
   # Password: autolmk_dev_password
   ```

4. **Start Server:**
   ```bash
   make run
   # OR
   make dev  # with hot reload
   # OR
   go run cmd/api/main.go
   ```

Server will start on: `http://localhost:8080`

---

## Testing Endpoints

### Health Check

```bash
curl http://localhost:8080/health
# Expected: OK
```

### API Info

```bash
curl http://localhost:8080/api/
# Expected: {"message":"Auto LMK API","version":"1.0.0"}
```

---

## Root Admin Endpoints

These endpoints do NOT require tenant domain.

### 1. Create Tenant

```bash
curl -X POST http://localhost:8080/api/admin/tenants \
  -H "Content-Type: application/json" \
  -d '{
    "domain": "test-showroom.localhost",
    "name": "Test Showroom",
    "whatsapp_number": "628123456789"
  }'
```

**Expected Response:**
```json
{
  "id": 1,
  "domain": "test-showroom.localhost",
  "name": "Test Showroom",
  "whatsapp_number": "628123456789",
  "pairing_status": "unpaired",
  "status": "active",
  "created_at": "2025-11-14T...",
  "updated_at": "2025-11-14T..."
}
```

### 2. List All Tenants

```bash
curl http://localhost:8080/api/admin/tenants
```

**Expected Response:**
```json
{
  "data": [
    {
      "id": 1,
      "domain": "test-showroom.localhost",
      "name": "Test Showroom",
      ...
    }
  ],
  "count": 1
}
```

### 3. Get Tenant by ID

```bash
curl http://localhost:8080/api/admin/tenants/1
```

---

## Tenant-Scoped Endpoints

These endpoints REQUIRE tenant domain in Host header for tenant isolation.

**Important:** Use `-H "Host: {tenant-domain}"` to specify tenant.

### 1. Create Car

```bash
curl -X POST http://localhost:8080/api/cars \
  -H "Host: showroom-jaya.localhost" \
  -H "Content-Type: application/json" \
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
    "description": "Toyota Avanza 2023 kondisi sangat baik",
    "is_featured": true
  }'
```

**Expected Response:**
```json
{
  "id": 1,
  "tenant_id": 1,
  "brand": "Toyota",
  "model": "Avanza",
  "year": 2023,
  "price": 250000000,
  ...
}
```

### 2. List Cars (with filters)

```bash
# All cars for tenant
curl -H "Host: showroom-jaya.localhost" \
  http://localhost:8080/api/cars

# Filter by brand
curl -H "Host: showroom-jaya.localhost" \
  "http://localhost:8080/api/cars?brand=Toyota"

# Filter by max price
curl -H "Host: showroom-jaya.localhost" \
  "http://localhost:8080/api/cars?max_price=300000000"

# Filter by transmission
curl -H "Host: showroom-jaya.localhost" \
  "http://localhost:8080/api/cars?transmission=automatic"

# Multiple filters
curl -H "Host: showroom-jaya.localhost" \
  "http://localhost:8080/api/cars?brand=Toyota&max_price=300000000&transmission=automatic"
```

### 3. Search Cars

```bash
curl -H "Host: showroom-jaya.localhost" \
  "http://localhost:8080/api/cars/search?q=Toyota"

curl -H "Host: showroom-jaya.localhost" \
  "http://localhost:8080/api/cars/search?q=matic"
```

### 4. Get Car by ID

```bash
curl -H "Host: showroom-jaya.localhost" \
  http://localhost:8080/api/cars/1
```

### 5. Update Car

```bash
curl -X PUT http://localhost:8080/api/cars/1 \
  -H "Host: showroom-jaya.localhost" \
  -H "Content-Type: application/json" \
  -d '{
    "price": 245000000,
    "description": "Toyota Avanza 2023 - Harga turun!"
  }'
```

### 6. Delete Car

```bash
curl -X DELETE http://localhost:8080/api/cars/1 \
  -H "Host: showroom-jaya.localhost"
```

---

## Testing Multi-Tenant Isolation

### Scenario: Verify tenant cannot access other tenant's data

1. **Create car for Tenant 1:**
```bash
curl -X POST http://localhost:8080/api/cars \
  -H "Host: showroom-jaya.localhost" \
  -H "Content-Type: application/json" \
  -d '{"brand":"Toyota","model":"Avanza","year":2023,"price":250000000}'
# Note the returned car ID
```

2. **Try to access from Tenant 2:**
```bash
curl -H "Host: mobilindo.localhost" \
  http://localhost:8080/api/cars/1
# Expected: 404 Not Found (car not found for this tenant)
```

3. **List cars for each tenant:**
```bash
# Tenant 1 should see their car
curl -H "Host: showroom-jaya.localhost" http://localhost:8080/api/cars

# Tenant 2 should see empty list
curl -H "Host: mobilindo.localhost" http://localhost:8080/api/cars
```

---

## Testing with Sample Data

If you ran the seed script, test with pre-populated data:

### List cars for Showroom Jaya:
```bash
curl -H "Host: showroom-jaya.localhost" http://localhost:8080/api/cars
# Should return 5 cars
```

### List cars for Mobilindo:
```bash
curl -H "Host: mobilindo.localhost" http://localhost:8080/api/cars
# Should return 3 luxury cars
```

### Search for Toyota:
```bash
curl -H "Host: showroom-jaya.localhost" \
  "http://localhost:8080/api/cars/search?q=Toyota"
# Should return Avanza and Fortuner
```

### Filter by budget:
```bash
curl -H "Host: showroom-jaya.localhost" \
  "http://localhost:8080/api/cars?max_price=300000000"
# Should return Avanza, Xpander, Terios
```

---

## Common Issues

### 1. Tenant Not Found
**Error:** `404: Tenant not found`
**Solution:**
- Check if domain exists in database
- Verify Host header matches tenant domain
- Use root admin endpoint to create tenant first

### 2. Database Connection Error
**Error:** `failed to connect to database`
**Solution:**
- Ensure Docker PostgreSQL is running: `docker-compose up -d`
- Check connection string in `.env`
- Test with: `psql -U autolmk -h localhost -d autolmk`

### 3. Migration Not Run
**Error:** `relation "tenants" does not exist`
**Solution:** Run migrations: `make migrate-up`

---

## Next Steps

After basic API testing works:

1. **Week 2-3:** Add authentication and session management
2. **Week 4-6:** Integrate WhatsApp bot and LLM
3. **Week 7-9:** Build web interface with HTMX

---

## Performance Testing

### Simple Load Test (using Apache Bench):

```bash
# Install ab (Apache Bench)
# macOS: brew install httpd
# Ubuntu: apt-get install apache2-utils

# Test health endpoint
ab -n 1000 -c 10 http://localhost:8080/health

# Test car listing
ab -n 100 -c 5 -H "Host: showroom-jaya.localhost" http://localhost:8080/api/cars
```

**Expected:**
- Health check: > 1000 req/sec
- Car listing: > 100 req/sec (with DB)

---

## Debugging

### Enable Debug Logging:

Set ENV=development in `.env`:
```env
ENV=development
```

This will show detailed SQL queries and request logs.

### Check Database:

```bash
psql -U autolmk -d autolmk -h localhost

# List tenants
SELECT * FROM tenants;

# List cars with tenant info
SELECT t.name, c.brand, c.model, c.price
FROM cars c
JOIN tenants t ON c.tenant_id = t.id;

# Check tenant isolation
SELECT tenant_id, COUNT(*) FROM cars GROUP BY tenant_id;
```

---

**Happy Testing! 🚀**
