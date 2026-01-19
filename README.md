[English](README.md) | [中文](README-CN.md)

# iComment - Minimalist Comment System

A lightweight, single-binary comment system designed for blogs.

## ✨ Features

- **Minimalist**: Single executable, SQLite database, no complex dependencies.
- **Docker Ready**: First-class Docker support with automated builds.
- **Secure**: Privacy-focused with comment approval workflow.
- **Easy Integration**: Drop-in HTML/JS snippet for any website.
- **Notifications**: Optional Bark (iOS) push notifications.

---

## 📖 User Guide

### Quick Start (Docker) - Recommended

Use the provided `docker-compose.yml` to get started instantly.

```yaml
services:
  icomment:
    image: claws059/icomment:latest
    restart: unless-stopped
    ports:
      - "7001:7001" # Public API
      - "7002:7002" # Admin UI (Protect this!)
    volumes:
      - ./data:/data
    environment:
      - ICOMMENT_BARK=your_bark_key # Optional
```

Start the service:

```bash
docker-compose up -d
```

### Configuration

You can configure iComment via Environment Variables (Docker) or Command-line Flags.

| Environment Var | Flag | Default | Description |
|-----------------|------|---------|-------------|
| `ICOMMENT_DB` | `-db` | `./comments.db` | Path to SQLite database |
| `ICOMMENT_PORT` | `-port` | `7001` | Public API & Static files port |
| `ICOMMENT_ADMIN_PORT` | `-admin-port` | `7002` | Admin Console port |
| `ICOMMENT_BARK` | `-bark` | (empty) | Bark Device Key for notifications |

### Integration

Add the following snippet to your blog templates:

```html
<div id="icomment"></div>
<script src="http://your-domain.com:7001/static/comment.js" 
        data-api="http://your-domain.com:7001"></script>
```

### Administration & Security

Access the Admin Console at `http://your-domain.com:7002/`.

> [!WARNING]
> **SECURITY NOTICE**: The Admin Console (Port 7002) listens on **all interfaces (0.0.0.0)**.
> You **MUST** protect this port in production using a firewall or reverse proxy (e.g., Nginx with Basic Auth). **Do not expose it publicly without protection.**

**Nginx Protection Example:**

```nginx
location /admin/ {
    proxy_pass http://127.0.0.1:7002/;
    # Enable Basic Auth
    auth_basic "Restricted Admin Area";
    auth_basic_user_file /etc/nginx/.htpasswd;
}
```

---

## 🛠 Developer Guide

### Build from Source

```bash
# Data directory
mkdir data

# Build
go build -o icomment .

# Run
./icomment
```

### Project Structure

```text
.
├── main.go              # Entry point & Server setup
├── config.go            # Config loading (Env/Flags)
├── handler.go           # HTTP Request Handlers
├── Dockerfile           # Multi-stage Docker build
├── .github/workflows/   # CI/CD (Docker Publish)
├── model/               # Data Structures
├── dao/                 # Database Access Layer
├── static/              # Frontend Assets (JS)
└── templates/           # Admin Admin UI Templates
```

### API Reference

**Public Endpoints (Port 7001)**
- `GET /api/comments?article_url=...` - List comments
- `POST /api/comments` - Submit comment

**Admin Endpoints (Port 7002)**
- `GET /comments` - List/Filter comments
- `PATCH /comments/:id/approve` - Approve comment
- `DELETE /comments/:id` - Delete comment

### Database Schema

iComment uses SQLite.

```sql
CREATE TABLE comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    article_url TEXT NOT NULL,
    nickname TEXT NOT NULL,
    email TEXT,
    content TEXT NOT NULL,
    status TEXT DEFAULT 'pending', -- pending/approved
    ...
);
```

---

## License

MIT
