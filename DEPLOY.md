# Nexus-MM Deployment Guide

## Prerequisites

- Ubuntu 22.04+ or Debian 12+
- Go 1.22+
- Node.js 18+
- Docker & Docker Compose
- Nginx
- A domain name (optional, for SSL)

## 1. Docker Compose Setup

```bash
# Clone the repository
git clone https://github.com/yujiawei/nexus-mm.git
cd nexus-mm

# Start all services
docker compose up -d
```

This starts:
- PostgreSQL (port 5432)
- Redis (port 6379)
- WuKongIM (ports 5001, 5100, 5200, 5300)
- MeiliSearch (port 7700)

## 2. Database Initialization

```bash
psql -h localhost -U nexus -d nexus_mm -f migrations/001_init.sql
```

## 3. Configuration

```bash
cp configs/nexus.yaml.example configs/nexus.yaml
```

Edit `configs/nexus.yaml`:

```yaml
server:
  host: "0.0.0.0"
  port: 9876

database:
  host: "localhost"
  port: 5432
  user: "nexus"
  password: "your-secure-password"
  dbname: "nexus_mm"
  sslmode: "disable"

jwt:
  secret: "your-long-random-secret-at-least-32-chars"
  expire_hour: 72

wukong:
  api_url: "http://localhost:5001"
  manager_token: "your-wukongim-manager-token"
  webhook_addr: "0.0.0.0:6979"

meilisearch:
  url: "http://localhost:7700"
  api_key: ""

redis:
  addr: "localhost:6379"
```

## 4. Build

```bash
# Build backend
go build -o nexus-mm ./cmd/server/

# Build frontend
cd web
npm install
npm run build
cd ..
```

## 5. Systemd Service

Create `/etc/systemd/system/nexus-mm.service`:

```ini
[Unit]
Description=Nexus-MM API Server
After=network.target postgresql.service redis.service

[Service]
Type=simple
User=nexus
Group=nexus
WorkingDirectory=/opt/nexus-mm
ExecStart=/opt/nexus-mm/nexus-mm -config /opt/nexus-mm/configs/nexus.yaml
Restart=always
RestartSec=5
LimitNOFILE=65535

Environment=GIN_MODE=release

[Install]
WantedBy=multi-user.target
```

```bash
# Install
sudo cp nexus-mm /opt/nexus-mm/
sudo cp -r configs /opt/nexus-mm/
sudo cp -r web/dist /opt/nexus-mm/web/

# Enable and start
sudo systemctl daemon-reload
sudo systemctl enable nexus-mm
sudo systemctl start nexus-mm
sudo systemctl status nexus-mm
```

## 6. Nginx Reverse Proxy

Create `/etc/nginx/sites-available/nexus-mm`:

```nginx
server {
    listen 80;
    server_name your-domain.com;

    # Frontend static files
    location / {
        root /opt/nexus-mm/web/dist;
        try_files $uri $uri/ /index.html;
    }

    # API proxy
    location /api/ {
        proxy_pass http://127.0.0.1:9876;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # Bot API proxy
    location /bot/ {
        proxy_pass http://127.0.0.1:9876;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    # Health check
    location /health {
        proxy_pass http://127.0.0.1:9876;
    }

    # skill.md endpoint
    location /skill.md {
        proxy_pass http://127.0.0.1:9876;
    }
}
```

```bash
sudo ln -s /etc/nginx/sites-available/nexus-mm /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

## 7. GCP Firewall Rules

```bash
# Allow HTTP
gcloud compute firewall-rules create nexus-mm-http \
  --allow tcp:80 \
  --target-tags=nexus-mm \
  --description="Allow HTTP for Nexus-MM"

# Allow HTTPS
gcloud compute firewall-rules create nexus-mm-https \
  --allow tcp:443 \
  --target-tags=nexus-mm \
  --description="Allow HTTPS for Nexus-MM"

# Allow API direct (if needed)
gcloud compute firewall-rules create nexus-mm-api \
  --allow tcp:9876 \
  --target-tags=nexus-mm \
  --description="Allow Nexus-MM API direct access"

# Allow WuKongIM WebSocket (for frontend WS connection)
gcloud compute firewall-rules create nexus-mm-ws \
  --allow tcp:5200 \
  --target-tags=nexus-mm \
  --description="Allow WuKongIM WebSocket"
```

## 8. SSL/TLS with Certbot

```bash
# Install certbot
sudo apt install -y certbot python3-certbot-nginx

# Obtain certificate
sudo certbot --nginx -d your-domain.com

# Auto-renewal (certbot sets up a systemd timer by default)
sudo certbot renew --dry-run
```

## 9. Monitoring

Check service status:

```bash
sudo systemctl status nexus-mm
sudo journalctl -u nexus-mm -f
```

Health check:

```bash
curl http://localhost:9876/health
```

## 10. Backup

```bash
# Database backup
pg_dump -h localhost -U nexus nexus_mm > backup_$(date +%Y%m%d).sql

# Restore
psql -h localhost -U nexus -d nexus_mm < backup_20260101.sql
```
