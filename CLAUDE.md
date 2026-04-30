# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Principles

### Solution Design Guidelines

- **Security & System Design First**: When providing solutions, always consider security implications, scalability, performance, and system design patterns. Think about potential vulnerabilities, data validation, error handling, and how the solution fits within the overall architecture.

- **User Experience (UX) Focus**: For frontend solutions, prioritize user experience. Consider loading states, error handling, accessibility, responsive design, intuitive interactions, and clear feedback to users. Solutions should feel natural and help users accomplish their goals efficiently.

## Common Development Commands

Recipes are defined in `justfile`. Run `just` with no args to list them.

### Backend (Go Server)

- **Start database**: `just postgres-up`
- **Stop database**: `just postgres-down`
- **Run migrations**: `just migrate-up`
- **Generate DB models**: `just gen-db-models` (after schema changes)
- **Build server**: `just server-build`
- **Run server**: `just server-run`
- **Run tests**: `just server-test`

### Frontend (SolidJS Web App)

- **Install dependencies**: `just web-install` (uses pnpm)
- **Development server**: `just web-dev`
- **Build production**: `just web-build`
- **Lint code**: `just web-lint`

### Full Application

- **Build both**: `just build`
- **Run complete app**: `just run` (builds both, serves web from Go server)
- **Run complete app in development mode**: `just dev` (builds both, serves web via vite dev server)
- **Run tests**: `just test`
- **Clean everything**: `just nuke`

## Architecture Overview

### Backend - Clean Architecture Pattern

**Domain Layer** (`/server/internal/domain/`):

- **auth/**: Authentication domain with JWT-based auth
- **media/**: File/folder management domain
- **profile/**: User profile management domain
- **sharing/**: File sharing functionality domain

Each domain follows CQRS pattern with:

- `commands.go` + `command_handlers.go`: Write operations
- `queries.go` + `query_handlers.go`: Read operations
- `*_sanitizer.go`: Input validation layer
- `repository.go`: Domain repository interface

**Infrastructure Layer** (`/server/internal/infrastructure/`):

- Repository implementations using Jet SQL generator
- Local file storage implementation
- JWT authentication infrastructure

**API Layer** (`/server/internal/api/`):

- Chi router with domain-specific API handlers
- JWT middleware for authentication
- DTOs for request/response models

**Workflows** (`/server/internal/workflows/`):

- Cross-domain operations (signup, signin)
- Handles transactional coordination between domains

### Frontend - SolidJS Architecture

**Component Structure**:

- `@sv/components/`: Reusable UI components
- `@sv/pages/`: Page-level components
- `@sv/store/`: Global state management
- `@sv/apis/`: API client layer
- `@sv/utils/`: Utility functions and constants

**Key Features**:

- Responsive design with Tailwind CSS
- Kobalte UI components for accessibility
- TanStack Query for server state management
- File upload with progress tracking and chunked uploads
- Real-time storage quota display with visual indicators

## Code Conventions

### Backend (Go)

- Follow conventions in @server/CONVENTIONS.md
- Use AppError for consistent error handling
- Document all possible app-errors in method comments
- Keep domain models free of infrastructure concerns
- Use transactions for cross-domain operations in workflows
- File naming: snake_case for files, single word for packages

### Frontend (TypeScript/SolidJS)

- Path aliases configured: `@sv/*` maps to `src/*`
- Use strongly-typed constants (see `@sv/utils/consts`)
- Follow SolidJS patterns for reactivity
- Prefer signals for primitive values and stores for complex state
- Component props use interfaces
- **Styling**: Prefer existing CSS classes from @web/src/index.css over direct Tailwind classes
  - Use semantic classes like `text-primary`, `btn btn-primary`, `border-border`, `bg-bg-muted`
  - Avoid hardcoded colors like `text-blue-500`, `border-gray-300`
  - Create new semantic classes in @web/src/index.css for common patterns

## Database

- PostgreSQL with Docker Compose
- Migrations in `/server/internal/infrastructure/internal/repository/internal/migrations/`
- Jet SQL generator for type-safe queries
- Models auto-generated in `/server/internal/infrastructure/internal/repository/internal/gen_jet/`

## File Storage

- Local file system storage implementation
- Configurable storage directory via environment variables
- Per-user storage quotas
- Storage usage tracking and display in UI
- Upload sessions track chunked uploads server-side with upfront quota reservation; foundation for future resumable uploads (end-to-end resume across client interruptions not yet wired)
- Chunked upload support for large files with concurrent-upload protection
- File categorization (image, video, audio, text, other)

## Testing

- Backend: Standard Go testing with testify
- Integration tests in `/server/internal/tests/integration/`
- No frontend tests currently configured

## Environment Configuration

- Backend configs:
  - Example config file: @server/.env.example
  - Development config file: @server/dev.env
  - Test config file: @server/test.env
  - Finally configs are loaded from @server/pkg/appconfig/config.go
- Uses strongly-typed configuration structs
- Database connection via environment variables
- Storage paths and JWT secrets configurable via environment variables
