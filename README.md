# SkyVault

<div align="center">

<!-- ![SkyVault Logo](logo.png) -->

**A self-hosted cloud storage solution for your files**

[![License: AGPLv3](https://img.shields.io/badge/License-AGPLv3-blue.svg)](LICENSE)
[![Version](https://img.shields.io/badge/version-1.0.0-green.svg)](https://github.com/vish9812/skyvault/releases)
[![Go Version](https://img.shields.io/badge/go-1.23-blue.svg)](https://golang.org)
[![SolidJS](https://img.shields.io/badge/SolidJS-1.9-blue.svg)](https://solidjs.com)

</div>

## 📋 Overview

SkyVault is a self-hosted cloud storage solution designed to help you securely store, organize, and share your files. It features a responsive mobile-first web UI and provides full control over your data.

### ✨ Features

- 🔐 **Secure Authentication**: JWT-based authentication system
- 📁 **Folder Management**: Create and navigate through folder structures
- 📤 **File Upload**: Upload files with support for resumable chunked uploads for large files
- 📥 **File Download**: Download your files anytime
- 💾 **Storage Quotas**: Per-user storage limits with real-time usage tracking
- 📊 **Storage Monitoring**: Visual storage quota display with usage indicators
- 📱 **Mobile-First UI**: Responsive design optimized for mobile devices
- 🎨 **Modern Interface**: Built with SolidJS and Tailwind CSS
- 🚀 **High Performance**: Go backend with clean architecture
- 🐳 **Easy Deployment**: Docker-based deployment with PostgreSQL

## 🚀 Quick Start

### Prerequisites

- Docker and Docker Compose installed
- At least 2GB of available disk space

### Installation via Docker

1. **Download the configuration file**

   ```bash
   wget https://raw.githubusercontent.com/vish9812/skyvault/main/.env.example -O .env
   ```

2. **Configure your environment**

   Edit the `.env` file and update the following critical values:

   ```bash
   # IMPORTANT: Change these to secure values!
   DB__PASS=your-secure-database-password
   AUTH__JWT__KEY=your-secure-jwt-key-at-least-32-characters-long
   ```

   Generate a secure JWT key:

   ```bash
   openssl rand -base64 32
   ```

3. **Download the Docker Compose file**

   ```bash
   wget https://raw.githubusercontent.com/vish9812/skyvault/main/docker-compose.prod.yml -O compose.yml
   ```

4. **Start SkyVault**

   ```bash
   docker compose up -d
   ```

5. **Access the application**

   Open your browser and navigate to:

   ```
   http://localhost:8090
   ```

   Or replace `localhost` with your server's IP address.

6. **Create your account**

   On first launch, create your user account through the signup page.

## 🔧 Configuration

### Environment Variables

SkyVault is configured through environment variables in the `.env` file:

| Variable                       | Description                           | Default           |
| ------------------------------ | ------------------------------------- | ----------------- |
| `SERVER__PORT`                 | Port to expose the application        | `8090`            |
| `DB__NAME`                     | PostgreSQL database name              | `skyvault`        |
| `DB__USER`                     | PostgreSQL username                   | `skyvault`        |
| `DB__PASS`                     | PostgreSQL password                   | ⚠️ **Required**   |
| `AUTH__JWT__KEY`               | JWT secret key (min 32 chars)         | ⚠️ **Required**   |
| `AUTH__JWT__TOKEN_TIMEOUT_MIN` | Token expiration in minutes           | `43200` (30 days) |
| `STORAGE__DEFAULT_QUOTA_MB`    | Default storage quota for new users   | `10240` (10GB)    |
| `LOG__LEVEL`                   | Logging level (debug/info/warn/error) | `info`            |

### Storage Quotas

SkyVault supports per-user storage quotas. You can configure the default quota for new users:

```bash
# Set default storage quota to 50GB for new users
STORAGE__DEFAULT_QUOTA_MB=51200

# Set to 100GB for generous storage
STORAGE__DEFAULT_QUOTA_MB=102400

# Set to 1TB for unlimited storage feel
STORAGE__DEFAULT_QUOTA_MB=1048576
```

**Features:**
- Each user has their own storage quota
- Real-time storage usage tracking displayed in the UI
- Visual indicators when approaching storage limits (80% = warning, 95% = critical)
- Automatic validation prevents uploads exceeding available storage
- Resumable chunked uploads for large files with concurrent upload protection

## 🛠️ Management

### Viewing Logs

```bash
# All logs
docker compose logs -f

# Application logs only
docker compose logs -f app

# Database logs only
docker compose logs -f db
```

### Updating SkyVault

```bash
# Pull the latest image
docker compose pull

# Restart with the new image
docker compose up -d
```

### Backup

Your data is stored in Docker volumes. To backup:

```bash
# Backup database
docker compose exec db pg_dump -U skyvault skyvault > backup.sql

# Backup uploaded files
docker run --rm -v skyvault_app-data:/data -v $(pwd):/backup alpine tar czf /backup/files-backup.tar.gz /data
```

### Restore

```bash
# Restore database
cat backup.sql | docker compose exec -T db psql -U skyvault skyvault

# Restore uploaded files
docker run --rm -v skyvault_app-data:/data -v $(pwd):/backup alpine tar xzf /backup/files-backup.tar.gz -C /
```

### Stopping SkyVault

```bash
# Stop services
docker compose down

# Stop and remove all data (⚠️ WARNING: This deletes everything!)
docker compose down -v
```

## 👩‍💻 Development

### Prerequisites

- Go 1.23 or higher
- Node.js 20 or higher
- pnpm
- PostgreSQL 16
- Task (taskfile.dev)

### Setup

1. **Clone the repository**

   ```bash
   git clone https://github.com/vish9812/skyvault.git
   cd skyvault
   ```

2. **Configure environment**

   ```bash
   cp server/dev.env.example server/dev.env
   # Edit server/dev.env with your database credentials
   ```

3. **Start the database**

   ```bash
   task postgres-up
   ```

4. **Start development servers**

   ```bash
   # Terminal 1: Start backend
   task server:run

   # Terminal 2: Start frontend
   task web:dev
   ```

   Access the app at `http://localhost:5173` (Vite dev server)

### Common Development Commands

```bash
# Build everything
task build

# Run all tests
task test

# Run server tests
task server:test

# Lint web code
task web:lint

# Generate DB models after schema changes
task gen-db-models

# Create a new migration
MIGRATION_FILE_NAME=add_something task migrate-create

# Clean everything
task nuke
```

### Project Structure

```
skyvault/
├── server/                 # Go backend
│   ├── cmd/               # Application entrypoint
│   ├── internal/
│   │   ├── domain/       # Domain layer (CQRS pattern)
│   │   ├── infrastructure/ # Infrastructure implementations
│   │   ├── api/          # HTTP API handlers
│   │   └── workflows/    # Cross-domain operations
│   └── pkg/              # Shared packages
├── web/                   # SolidJS frontend
│   └── src/
│       ├── components/   # UI components
│       ├── pages/        # Page components
│       ├── store/        # State management
│       └── apis/         # API client
└── Taskfile.yml          # Task automation
```

## 🏗️ Architecture

### Backend

- **Language**: Go 1.23
- **Architecture**: Clean Architecture with CQRS pattern
- **Database**: PostgreSQL 16
- **Authentication**: JWT tokens
- **API**: RESTful HTTP API with Chi router

### Frontend

- **Framework**: SolidJS 1.9
- **Styling**: Tailwind CSS 4
- **UI Components**: Kobalte
- **State Management**: Solid Signals + TanStack Query
- **Build Tool**: Vite

### Storage

- **Type**: Local filesystem storage
- **Features**:
  - Per-user storage quotas
  - Real-time usage tracking
  - Resumable chunked uploads for large files
  - Streaming downloads
  - Concurrent upload protection
- **Configuration**: Default quota configurable via `STORAGE__DEFAULT_QUOTA_MB`

## 🗺️ Roadmap

### Completed ✅

- JWT-based authentication
- Folder creation and navigation
- File upload with resumable chunked uploads
- File download
- Per-user storage quotas with real-time tracking
- Storage usage visualization in UI

### In Progress 🚧

- Epic 1: File Operations (rename, move, delete)
- Epic 2: Folder Operations (rename, move, delete)

### Planned 📋

- Epic 3: Contact Management System
- Epic 4: Core File Sharing
- Epic 5: Shared Content Management
- Epic 6: Advanced Sharing Features

See [TODO.md](TODO.md) for detailed roadmap.

## 🤝 Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines

- Follow the conventions in [CONVENTIONS.md](server/CONVENTIONS.md)
- Write tests for new features
- Update documentation as needed
- Ensure all tests pass before submitting PR

## 📄 License

This project is licensed under the GNU AGPLv3 License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [SolidJS](https://solidjs.com) - Reactive UI framework
- [Go](https://golang.org) - Backend language
- [Chi](https://github.com/go-chi/chi) - HTTP router
- [Kobalte](https://kobalte.dev) - Accessible UI components
- [Tailwind CSS](https://tailwindcss.com) - Utility-first CSS

## 📞 Support

- **Issues**: [GitHub Issues](https://github.com/vish9812/skyvault/issues)
- **Discussions**: [GitHub Discussions](https://github.com/vish9812/skyvault/discussions)

## ⚠️ Security

If you discover a security vulnerability, please email [vishapps@outlook.com](mailto:vishapps@outlook.com) instead of using the issue tracker.

---

<div align="center">

Made with ❤️ by [Vish](https://github.com/vish9812)

⭐ Star this repository if you find it useful!

</div>
