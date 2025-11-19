# Changelog

All notable changes to SkyVault will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2025-10-02

### 🎉 Initial Release

This is the first public release of SkyVault, a self-hosted cloud storage solution with a modern web interface.

### ✨ Features

#### Authentication
- JWT-based authentication system
- Secure password hashing with bcrypt
- Token-based session management with configurable expiration
- Sign up and sign in functionality

#### File Management
- File upload with drag-and-drop support
- Chunked upload for large files (up to 10GB)
  - Automatic chunking for files over 5GB
  - Configurable chunk size (default: 100MB)
- File download with streaming support
- File categorization (image, video, audio, text, other)
- File size limits configurable via environment variables

#### Folder Management
- Create nested folder structures
- Navigate through folder hierarchy
- Breadcrumb navigation
- Parent folder navigation

#### User Interface
- Mobile-first responsive design
- Built with SolidJS for high performance
- Styled with Tailwind CSS 4
- Accessible UI components with Kobalte
- Modern gradient design
- Real-time upload progress tracking
- Intuitive file and folder interactions

#### Backend
- Clean Architecture with CQRS pattern
- Domain-driven design
- RESTful API with Chi router
- PostgreSQL database with migrations
- Local filesystem storage
- Comprehensive error handling
- Structured logging with zerolog
- Health check endpoint for monitoring

#### Deployment
- Docker-based deployment
- Multi-stage Dockerfile for optimized image size
- Docker Compose configuration for production
- Automatic database migrations on startup
- Volume-based data persistence
- Health checks for both database and application
- Support for both amd64 and arm64 architectures

#### Developer Experience
- Task-based build system with Taskfile
- Separate development and test environments
- Hot reload for frontend development
- Comprehensive development documentation
- Code conventions and guidelines
- Integration tests

### 🏗️ Architecture

#### Backend Stack
- **Language**: Go 1.23
- **Framework**: Chi router
- **Database**: PostgreSQL 16
- **ORM**: Jet (type-safe SQL)
- **Authentication**: golang-jwt/jwt
- **Logging**: zerolog

#### Frontend Stack
- **Framework**: SolidJS 1.9
- **Build Tool**: Vite 6
- **Styling**: Tailwind CSS 4
- **UI Components**: Kobalte 0.13
- **State Management**: TanStack Query 5
- **Router**: @solidjs/router

### 📦 Docker Images

- Multi-platform support: `linux/amd64`, `linux/arm64`
- Available on GitHub Container Registry
- Optimized multi-stage builds
- Alpine-based runtime (~50MB compressed)

### 🔐 Security

- Secure password hashing with bcrypt
- JWT tokens with configurable expiration
- Database connection over SSL (configurable)
- Input validation and sanitization
- Error context without sensitive data leakage

### 📝 Documentation

- Comprehensive README with installation instructions
- Docker deployment guide
- Development setup instructions
- Configuration reference
- API documentation through code
- Release guide

### 🐛 Known Limitations

- Single user support (multi-user planned for future)
- No file sharing functionality yet (planned for future releases)
- No file/folder rename, move, or delete operations (in progress)
- Local storage only (cloud storage planned for future)

### 🔮 Coming Soon

See [TODO.md](TODO.md) for planned features:
- File operations (rename, move, delete)
- Folder operations (rename, move, delete)
- Contact management system
- File sharing with password protection
- Shared content management
- Advanced sharing features

---

## Release Format

### [Version] - YYYY-MM-DD

#### Added
- New features

#### Changed
- Changes to existing functionality

#### Deprecated
- Features that will be removed in upcoming releases

#### Removed
- Removed features

#### Fixed
- Bug fixes

#### Security
- Security improvements or fixes

---

[1.0.0]: https://github.com/yourusername/skyvault/releases/tag/v1.0.0
