# Deployment Guide - Auto LMK Platform

**Version:** 1.0.0
**Last Updated:** 2025-11-14

## 🎯 Deployment Overview

This guide covers deploying the Auto LMK platform to production.

---

## 📋 Prerequisites

### Required
- Linux server (Ubuntu 22.04 LTS recommended)
- 2GB RAM minimum (4GB recommended)
- 20GB disk space
- Domain name
- SSL certificate capability
- PostgreSQL 15+
- Go 1.21+

### Optional but Recommended
- Docker & Docker Compose
- Nginx (reverse proxy)
- Let's Encrypt (free SSL)
- Systemd (service management)

---

## 🚀 Deployment Options

### Option 1: Docker Compose (Recommended)

**Pros:** Easiest, reproducible, isolated
**Use Case:** Small to medium deployments

**Steps:**

1. **Clone Repository:**
```bash
git clone https://github.com/yourusername/auto-lmk.git
cd auto-lmk
```

2. **Configure Environment:**
```bash
cp .env.example .env
vim .env
```

Edit `.env`:
```env
# Production settings
ENV=production
PORT=8080

# Database
DB_HOST=postgres
DB_PORT=5432
DB_USER=autolmk
DB_PASSWORD=YOUR_SECURE_PASSWORD_HERE
DB_NAME=autolmk
DB_SSLMODE=require

# LLM Provider
LLM_PROVIDER=openai
LLM_API_KEY=YOUR_OPENAI_API_KEY
LLM_MODEL=gpt-4o-mini

# WhatsApp
WHATSAPP_SESSION_PATH=./whatsapp_sessions

# Security
JWT_SECRET=YOUR_RANDOM_JWT_SECRET_HERE
```

3. **Create Production Docker Compose:**
```yaml
# docker-compose.prod.yml
version: '3.8'

services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - ENV=production
    env_file:
      - .env
    depends_on:
      - postgres
    restart: unless-stopped
    volumes:
      - ./uploads:/app/uploads
      - ./whatsapp_sessions:/app/whatsapp_sessions

  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_USER: ${DB_USER}
      POSTGRES_PASSWORD: ${DB_PASSWORD}
      POSTGRES_DB: ${DB_NAME}
    volumes:
      - postgres_data:/var/lib/postgresql/data
    restart: unless-stopped

volumes:
  postgres_data:
```

4. **Build and Deploy:**
```bash
docker-compose -f docker-compose.prod.yml up -d --build
```

5. **Run Migrations:**
```bash
docker-compose -f docker-compose.prod.yml exec app migrate -path /app/migrations -database "${DATABASE_URL}" up
```

---

### Option 2: Systemd Service (Traditional)

**Pros:** Direct control, no Docker overhead
**Use Case:** VPS deployment

**Steps:**

1. **Build Binary:**
```bash
CGO_ENABLED=0 GOOS=linux go build -o auto-lmk cmd/api/main.go
```

2. **Copy to Server:**
```bash
scp auto-lmk user@server:/opt/auto-lmk/
scp -r migrations user@server:/opt/auto-lmk/
```

3. **Create Systemd Service:**
```bash
sudo vim /etc/systemd/system/auto-lmk.service
```

```ini
[Unit]
Description=Auto LMK Car Sales Platform
After=network.target postgresql.service

[Service]
Type=simple
User=auto-lmk
WorkingDirectory=/opt/auto-lmk
EnvironmentFile=/opt/auto-lmk/.env
ExecStart=/opt/auto-lmk/auto-lmk
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
```

4. **Enable and Start:**
```bash
sudo systemctl daemon-reload
sudo systemctl enable auto-lmk
sudo systemctl start auto-lmk
sudo systemctl status auto-lmk
```

---

## 🌐 Nginx Reverse Proxy

**Purpose:** SSL termination, load balancing, static file serving

**Configuration:**

```nginx
# /etc/nginx/sites-available/auto-lmk

# Redirect HTTP to HTTPS
server {
    listen 80;
    server_name showroom-jaya.com www.showroom-jaya.com;
    return 301 https://$server_name$request_uri;
}

# Main HTTPS server
server {
    listen 443 ssl http2;
    server_name showroom-jaya.com www.showroom-jaya.com;

    # SSL Configuration
    ssl_certificate /etc/letsencrypt/live/showroom-jaya.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/showroom-jaya.com/privkey.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;

    # Proxy to Go application
    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # Static files (if serving from Nginx)
    location /uploads/ {
        alias /opt/auto-lmk/uploads/;
        expires 30d;
        add_header Cache-Control "public, immutable";
    }

    # Health check endpoint
    location /health {
        proxy_pass http://localhost:8080/health;
        access_log off;
    }
}

# Admin subdomain
server {
    listen 443 ssl http2;
    server_name admin.platform.com;

    ssl_certificate /etc/letsencrypt/live/platform.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/platform.com/privkey.pem;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

**Enable:**
```bash
sudo ln -s /etc/nginx/sites-available/auto-lmk /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

---

## 🔒 SSL Certificate (Let's Encrypt)

```bash
# Install Certbot
sudo apt install certbot python3-certbot-nginx

# Obtain certificate
sudo certbot --nginx -d showroom-jaya.com -d www.showroom-jaya.com

# Auto-renewal (already set up by certbot)
sudo systemctl status certbot.timer
```

---

## 🗄️ Database Setup

### PostgreSQL Installation

```bash
# Ubuntu
sudo apt update
sudo apt install postgresql-15 postgresql-contrib

# Start PostgreSQL
sudo systemctl start postgresql
sudo systemctl enable postgresql

# Create user and database
sudo -u postgres psql

CREATE USER autolmk WITH PASSWORD 'secure_password_here';
CREATE DATABASE autolmk OWNER autolmk;
GRANT ALL PRIVILEGES ON DATABASE autolmk TO autolmk;
\q
```

### Run Migrations

```bash
migrate -path ./migrations \
  -database "postgresql://autolmk:password@localhost:5432/autolmk?sslmode=disable" \
  up
```

### Seed Data (Optional)

```bash
psql -U autolmk -d autolmk -h localhost < scripts/seed.sql
```

---

## 📊 Monitoring

### Health Checks

```bash
# Simple uptime monitoring
*/5 * * * * curl -f http://localhost:8080/health || systemctl restart auto-lmk
```

### Logging

**Systemd logs:**
```bash
sudo journalctl -u auto-lmk -f
```

**Application logs:**
```bash
tail -f /var/log/auto-lmk/app.log
```

### Metrics (Optional)

Install Prometheus + Grafana for advanced monitoring.

---

## 🔄 Updates and Maintenance

### Deploy New Version

```bash
# Docker
git pull
docker-compose -f docker-compose.prod.yml build
docker-compose -f docker-compose.prod.yml up -d

# Systemd
git pull
go build -o auto-lmk cmd/api/main.go
sudo systemctl restart auto-lmk
```

### Database Backup

```bash
# Automated daily backup
0 2 * * * pg_dump -U autolmk autolmk | gzip > /backups/autolmk-$(date +\%Y\%m\%d).sql.gz
```

### Rollback

```bash
# Docker
docker-compose -f docker-compose.prod.yml down
git checkout <previous-commit>
docker-compose -f docker-compose.prod.yml up -d

# Database rollback
migrate -path ./migrations -database "${DATABASE_URL}" down 1
```

---

## 🔐 Security Checklist

- [ ] Change all default passwords
- [ ] Enable firewall (ufw or iptables)
- [ ] Configure SSH key-only access
- [ ] Set up fail2ban
- [ ] Enable SSL/TLS
- [ ] Configure rate limiting
- [ ] Regular security updates
- [ ] Database encryption at rest
- [ ] Environment variables (no hardcoded secrets)
- [ ] Regular backups
- [ ] Monitoring and alerting

---

## 🌍 Multi-Tenant Domain Setup

### DNS Configuration

For each tenant:

**Option A: Subdomain**
```
A Record: showroom-jaya.platform.com → YOUR_SERVER_IP
```

**Option B: Custom Domain**
```
A Record: showroom-jaya.com → YOUR_SERVER_IP
CNAME: www.showroom-jaya.com → showroom-jaya.com
```

### SSL for Multiple Domains

```bash
# Single certificate with multiple domains
sudo certbot --nginx \
  -d platform.com \
  -d admin.platform.com \
  -d showroom-jaya.platform.com \
  -d mobilindo.platform.com

# Or wildcard (requires DNS challenge)
sudo certbot certonly --dns-cloudflare \
  -d platform.com \
  -d *.platform.com
```

---

## 📈 Scaling

### Horizontal Scaling

```yaml
# docker-compose with replicas
services:
  app:
    deploy:
      replicas: 3
    # ... rest of config
```

### Load Balancer (Nginx)

```nginx
upstream auto_lmk_backend {
    least_conn;
    server localhost:8080;
    server localhost:8081;
    server localhost:8082;
}

server {
    location / {
        proxy_pass http://auto_lmk_backend;
    }
}
```

### Database Scaling

- Read replicas for queries
- Connection pooling (already configured)
- Indexes (already added)
- Query optimization

---

## 🚨 Troubleshooting

### Common Issues

**1. Database Connection Failed**
```bash
# Check PostgreSQL status
sudo systemctl status postgresql

# Check connection
psql -U autolmk -d autolmk -h localhost

# Check firewall
sudo ufw status
```

**2. Port Already in Use**
```bash
# Find process
sudo lsof -i :8080

# Kill process
sudo kill -9 <PID>
```

**3. Permission Issues**
```bash
# Fix ownership
sudo chown -R auto-lmk:auto-lmk /opt/auto-lmk
sudo chmod -R 755 /opt/auto-lmk
```

---

## 📞 Support

For deployment issues, check:
1. `docs/api-testing-guide.md`
2. Application logs: `journalctl -u auto-lmk -f`
3. Database logs: `/var/log/postgresql/`
4. Nginx logs: `/var/log/nginx/error.log`

---

**Deployment Status Checklist:**

- [ ] Server provisioned
- [ ] Domain configured
- [ ] SSL certificate installed
- [ ] Database created and migrated
- [ ] Application deployed
- [ ] Nginx configured
- [ ] Firewall configured
- [ ] Monitoring set up
- [ ] Backups automated
- [ ] Health checks enabled
- [ ] Documentation updated

---

**Ready for production! 🚀**
