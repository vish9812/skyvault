# Repository Guidelines

## Project Structure & Module Organization

SkyVault is split into a Go backend and a SolidJS frontend. Backend code lives in `server/`: `cmd/main.go` starts the app, `internal/api` handles HTTP, `internal/domain` contains business logic, `internal/infrastructure` contains persistence and migrations, and `pkg` holds shared utilities. Tests live beside packages as `*_test.go` or in `server/internal/tests/integration`.

Frontend code lives in `web/src`: `pages` for routes, `components` for reusable UI, `apis` for typed API clients, `store` for app context, `utils` for helpers, and `assets` for static files. Root files include Docker config and the `justfile` task runner.

Application progress and next tasks are maintained in `roadmap.md`; keep roadmap status updates there instead of duplicating them in other docs.

## Product & Engineering Principles

SkyVault is a file-storage application, so security, data integrity, and user trust are the default priorities. Treat authentication, authorization, path handling, quota enforcement, upload validation, error handling, and data ownership checks as core product behavior rather than incidental implementation details.

Design backend changes with the full system in mind: keep domain rules explicit, preserve clean boundaries, consider transactionality and rollback behavior, and avoid shortcuts that make storage state, database state, or quota accounting diverge. Prefer simple, auditable flows over clever abstractions.

Design frontend changes around a clear user experience. File operations should provide responsive feedback, accessible controls, useful loading and error states, and predictable behavior on slow or failed networks. UI changes should help users understand what happened, what is happening, and what they can do next.

## Architecture Notes

Backend domains under `server/internal/domain` follow a clean architecture style. Keep write operations in `commands.go` and `command_handlers.go`, read operations in `queries.go` and `query_handlers.go`, input validation in `*_sanitizer.go`, and repository interfaces in the domain layer. Use `server/internal/workflows` only for cross-domain write operations such as signup/signin coordination, and keep transactional orchestration there.

Infrastructure belongs in `server/internal/infrastructure`: repository implementations use Jet-generated SQL models, migrations live in `server/internal/infrastructure/internal/repository/internal/migrations`, and generated Jet models live in `server/internal/infrastructure/internal/repository/internal/gen_jet`. API handlers and DTOs belong in `server/internal/api`.

Frontend imports use the `@sv/*` alias for `web/src/*`. The app uses SolidJS, Tailwind CSS, Kobalte UI components, and TanStack Query; follow existing patterns for server state, upload progress, quota display, and accessible UI components.

## Build, Test, and Development Commands

Use `just` from the repository root.

- `just postgres-up` / `just postgres-down`: manage the development PostgreSQL container.
- `just migrate-up`: apply database migrations.
- `just gen-db-models`: regenerate Jet SQL models after schema changes.
- `just server-run`: build and run the Go server with `server/dev.env`.
- `just server-test`: run backend tests via `go test ./...`.
- `just web-dev`: install web dependencies with `pnpm` and start Vite.
- `just web-build`: install dependencies and build the frontend.
- `just web-lint`: run the frontend linter.
- `just dev`: print the two-terminal development flow.
- `just run`: build both apps and serve the web build from the Go server.
- `just build`: build both frontend and backend.
- `just test`: run backend tests via `go test ./...`.
- `just postgres-test-up`: start the integration-test database.

For frontend-only work, run from `web/`, for example `pnpm run build` or `pnpm run dev`. Use `pnpm`; `preinstall` rejects other package managers.

## Coding Style & Naming Conventions

Format Go with `gofmt`; keep packages lowercase and one word where possible. Per `server/CONVENTIONS.md`, name package-matching files after the package and use `snake_case.go` for other Go files. Keep domain models free of infrastructure tags, use API DTOs for request/response shapes, define domain repository interfaces in the domain layer, and wrap domain errors with location-aware `AppError`s via `NewAppError(err, "location")`.

Document possible app-errors in method comments using the `App Errors:` format. Use strongly typed configuration structs and enums where practical. For cross-domain writes, use transactions through repository transactional support and commit or roll back in the workflow layer.

Frontend files use TypeScript/TSX with existing SolidJS patterns. Prefer small components, typed API models, and existing utilities before adding abstractions. Use signals for primitive local state and stores for complex state when that matches existing code. Component props should be typed with interfaces.

Prefer semantic CSS classes from `web/src/index.css` over hardcoded Tailwind colors. Use classes such as `text-primary`, `btn btn-primary`, `border-border`, and `bg-bg-muted`; add reusable semantic classes there when a pattern repeats.

## Testing Guidelines

Add Go unit tests beside the package under test as `*_test.go`. Put broader database/API flows in `server/internal/tests/integration`. Start the test database with `just postgres-test-up`, run `just test`, then stop it with `just postgres-test-down`. There is no frontend test script currently, so validate UI changes with `pnpm run build`.

Use testify where it is already established. For database changes, add migrations, run them with `just migrate-up`, regenerate Jet models with `just gen-db-models`, and include tests for repository or API behavior when practical.

## Environment & Storage

Backend configuration starts from `server/.env.example`, `server/dev.env.example`, or `server/test.env.example`; local concrete files such as `server/dev.env` and `server/test.env` must not contain committed secrets. Configuration loading is centralized in `server/pkg/appconfig/config.go`.

Storage is local-filesystem based and configurable through environment variables. Preserve per-user quotas, storage usage accounting, upload-session reservation behavior, and concurrent-upload protections when changing upload or file-management code.

## Commit & Pull Request Guidelines

Recent history uses short imperative messages, often Conventional Commit prefixes such as `fix:` and `chore:`. Keep subjects concise, for example `fix: handle expired upload sessions`. Pull requests should describe the change, list verification commands, link issues, and include screenshots for visible UI changes.

## Security & Configuration Tips

Do not commit real secrets or local data. Start from `.env.example`, `server/dev.env.example`, or `server/test.env.example`, and keep JWT keys and database passwords local. Uploaded files and generated database data should remain outside version control.
