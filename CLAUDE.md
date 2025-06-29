# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Principles

### Solution Design Guidelines

- **Security & System Design First**: When providing solutions, always consider security implications, scalability, performance, and system design patterns. Think about potential vulnerabilities, data validation, error handling, and how the solution fits within the overall architecture.

- **User Experience (UX) Focus**: For frontend solutions, prioritize user experience. Consider loading states, error handling, accessibility, responsive design, intuitive interactions, and clear feedback to users. Solutions should feel natural and help users accomplish their goals efficiently.

## Common Development Commands

### Backend (Go Server)

- **Start database**: `task postgres-up`
- **Stop database**: `task postgres-down`
- **Run migrations**: `task migrate-up`
- **Generate DB models**: `task gen-db-models` (after schema changes)
- **Build server**: `task server:build`
- **Run server**: `task server:run`
- **Run tests**: `task server:test`

### Frontend (SolidJS Web App)

- **Install dependencies**: `task web:install` (uses pnpm)
- **Development server**: `task web:dev`
- **Build production**: `task web:build`
- **Lint code**: `task web:lint`

### Full Application

- **Build both**: `task build`
- **Run complete app**: `task run` (builds both, serves web from Go server)
- **Run complete app in development mode**: `task dev` (builds both, serves web via vite dev server)
- **Run tests**: `task test`
- **Clean everything**: `task nuke`

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
- Chunked upload support for large files
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
