# Repository Guidelines

## Project Structure & Module Organization

SkyVault is split into a Go backend and a SolidJS frontend. Backend code lives in `server/`: `cmd/main.go` starts the app, `internal/api` handles HTTP, `internal/domain` contains business logic, `internal/infrastructure` contains persistence and migrations, and `pkg` holds shared utilities. Tests live beside packages as `*_test.go` or in `server/internal/tests/integration`.

Frontend code lives in `web/src`: `pages` for routes, `components` for reusable UI, `apis` for typed API clients, `store` for app context, `utils` for helpers, and `assets` for static files. Root files include Docker config and the `justfile` task runner.

## Build, Test, and Development Commands

Use `just` from the repository root.

- `just postgres-up` / `just postgres-down`: manage the development PostgreSQL container.
- `just server-run`: build and run the Go server with `server/dev.env`.
- `just web-dev`: install web dependencies with `pnpm` and start Vite.
- `just dev`: print the two-terminal development flow.
- `just build`: build both frontend and backend.
- `just test`: run backend tests via `go test ./...`.
- `just postgres-test-up`: start the integration-test database.

For frontend-only work, run from `web/`, for example `pnpm run build` or `pnpm run dev`. Use `pnpm`; `preinstall` rejects other package managers.

## Coding Style & Naming Conventions

Format Go with `gofmt`; keep packages lowercase and one word where possible. Per `server/CONVENTIONS.md`, name package-matching files after the package and use `snake_case.go` for other Go files. Keep domain models free of infrastructure tags and wrap domain errors with location-aware `AppError`s.

Frontend files use TypeScript/TSX with existing SolidJS patterns. Prefer small components, typed API models, and existing utilities before adding abstractions.

## Testing Guidelines

Add Go unit tests beside the package under test as `*_test.go`. Put broader database/API flows in `server/internal/tests/integration`. Start the test database with `just postgres-test-up`, run `just test`, then stop it with `just postgres-test-down`. There is no frontend test script currently, so validate UI changes with `pnpm run build`.

## Commit & Pull Request Guidelines

Recent history uses short imperative messages, often Conventional Commit prefixes such as `fix:` and `chore:`. Keep subjects concise, for example `fix: handle expired upload sessions`. Pull requests should describe the change, list verification commands, link issues, and include screenshots for visible UI changes.

## Security & Configuration Tips

Do not commit real secrets or local data. Start from `.env.example`, `server/dev.env.example`, or `server/test.env.example`, and keep JWT keys and database passwords local. Uploaded files and generated database data should remain outside version control.
