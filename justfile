set dotenv-load
set dotenv-path := "server/dev.env"

DB_MIGRATION_PATH := "./internal/infrastructure/internal/repository/internal/migrations"
DB_GEN_MODELS_PATH := "./internal/infrastructure/internal/repository/internal/gen_jet"
SERVER_BUILD_OUTPUT := "./bin/server"
STATIC_DIR := "./static"
SERVER_STATIC_DIR := "./server/static"

# List available recipes.
default:
    @just --list

# ───────── Server ─────────

# Start a postgres container and create the DB.
[working-directory('server')]
postgres-up:
    mkdir -p $SERVER__DATA_DIR
    docker compose --env-file dev.env up -d

# Stop and remove the postgres container.
[working-directory('server')]
postgres-down:
    docker compose --env-file dev.env down

# Start a postgres container for integration tests on port 15433. Runs alongside the dev container.
# Sources test.env to override the dev vars just exports — docker compose's --env-file
# loses to existing shell vars, so we have to overwrite them ourselves.
[working-directory('server')]
postgres-test-up:
    #!/usr/bin/env bash
    set -euo pipefail
    set -a; source test.env; set +a
    mkdir -p "$SERVER__DATA_DIR/db"
    docker compose --env-file test.env -p skyvault-test up -d

# Stop and remove the test postgres container.
[working-directory('server')]
postgres-test-down:
    #!/usr/bin/env bash
    set -euo pipefail
    set -a; source test.env; set +a
    docker compose --env-file test.env -p skyvault-test down

# Create a new DB.
[working-directory('server')]
create-db:
    docker exec -it $DB__CONTAINER__NAME createdb --username=$DB__USER $DB__NAME

# Drop the DB.
[working-directory('server')]
drop-db:
    docker exec -it $DB__CONTAINER__NAME dropdb --username=$DB__USER $DB__NAME

# Migrate DB to the latest SQL file.
[working-directory('server')]
migrate-up:
    migrate -path {{ DB_MIGRATION_PATH }} -database $DB__DSN -verbose up

# Migrate DB down to 1 previous version.
[working-directory('server')]
migrate-down:
    migrate -path {{ DB_MIGRATION_PATH }} -database $DB__DSN -verbose down 1

# Create the next up and down sql files. Pass NAME=my_migration.
[working-directory('server')]
migrate-create NAME:
    migrate create -ext sql -dir {{ DB_MIGRATION_PATH }} -seq {{ NAME }}

# Generate DB models from the live schema.
[working-directory('server')]
gen-db-models:
    jet -dsn="$DB__DSN" -schema=public -path={{ DB_GEN_MODELS_PATH }}

# Build the server application.
[working-directory('server')]
server-build:
    mkdir -p bin
    go build -o {{ SERVER_BUILD_OUTPUT }} ./cmd/main.go

# Build and run the server application.
[working-directory('server')]
server-run: server-build
    {{ SERVER_BUILD_OUTPUT }} -dev -env dev.env

# Run the server tests.
[working-directory('server')]
server-test:
    go test ./...

# Remove containers and generated server data.
[working-directory('server')]
server-nuke: postgres-down
    sudo rm -rf $SERVER__DATA_DIR
    rm -rf bin
    rm -rf {{ STATIC_DIR }}

# ───────── Web ─────────

# Install web dependencies.
[working-directory('web')]
web-install:
    pnpm install

# Build the web application.
[working-directory('web')]
web-build: web-install
    pnpm run build

# Start the web development server.
[working-directory('web')]
web-dev: web-install
    pnpm run dev

# Run the web linter.
[working-directory('web')]
web-lint: web-install
    pnpm run lint

# Remove web generated files and dependencies.
[working-directory('web')]
web-nuke:
    rm -rf node_modules
    rm -rf dist

# ───────── Full Application ─────────

# Build both server and web applications.
build: web-build server-build

# Build and run the server (which serves the web app).
run: build
    rm -rf {{ SERVER_STATIC_DIR }}
    ln -s ../web/dist {{ SERVER_STATIC_DIR }}
    @just server-run

# Show how to run dev mode.
dev:
    @echo "Development mode requires 2 terminals:"
    @echo "  Terminal 1- just server-run"
    @echo "  Terminal 2- just web-dev"
    @echo ""
    @echo "Or use VS Code:"
    @echo "  1. Press F5 to start Go server"
    @echo "  2. Run 'just web-dev' in terminal"
    @echo ""
    @echo "Then open http://localhost:3000 in your browser"

# Run all tests.
test: server-test

# Clean all generated files.
nuke: web-nuke server-nuke
