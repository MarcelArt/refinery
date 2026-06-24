# AGENTS.md

Guidance for AI agents working in this repo.

## Hard rules

- **Do NOT run `serve`, `migrate`, `migrate-force`, `make dev`, `make go`, or any command that starts the server or touches the database.** These are for the human to run. Build/typecheck/compile edits only.
- **Do not cross the backend/frontend boundary in a single task.** If a task targets the backend (`internal/v1`, `internal/app`, `internal/entities`, etc.), do not also edit the frontend, and vice versa. If a change genuinely requires both, stop and tell the user instead of doing both yourself.
- **Never fabricate library/package APIs.** Before writing code that uses any library (Gin, GORM, HTMX, Alpine.js, dig, cobra, swaggo, etc.), look up the real API via the Context7 MCP (`resolve-library-id` → `query-docs`). Do not rely on memory for signatures, generics syntax, or HTMX attributes.
- **Retain the existing project structure and layering.** Do not introduce new top-level packages or reorganize without being asked.

## Backend / Frontend split

- **Backend** (REST API): `main.go`, `cmd/`, `internal/app/`, `internal/configs/`, `internal/container/`, `internal/entities/`, `internal/enums/`, `internal/common/`, `internal/v1/`, `pkg/`.
- **Frontend** (server-rendered UI): `internal/web/` — see its own section below.

## File naming conventions (important)

Files are split per domain, one route per file, but multiple related structs/functions may share a file when they belong to the same domain. Match existing suffixes exactly:

| Layer          | Directory                          | Suffix           | Example                |
|----------------|------------------------------------|------------------|------------------------|
| Entity         | `internal/entities/`               | `.entity.go`     | `user.entity.go`       |
| Model (DTO)    | `internal/v1/models/`              | `.model.go`      | `user.model.go`        |
| Repository     | `internal/v1/repositories/`        | `.repo.go`       | `user.repo.go`         |
| Service        | `internal/v1/services/`            | `.service.go`    | `user.service.go`      |
| Handler        | `internal/v1/handlers/`            | `.handler.go`    | `user.handler.go`      |
| Route          | `internal/v1/routes/`              | `.route.go`      | `user.route.go`        |

- A new domain (e.g. `foo`) means: `foo.entity.go`, `foo.model.go`, `foo.repo.go`, `foo.service.go`, `foo.handler.go`, `foo.route.go`. All six are usually needed.
- A file like `user.viewmodel.go` / `user.model.go` may hold multiple related types (e.g. `CreateUserViewModel`, `UserPageViewModel`). Group by domain, not one-type-per-file.
- Do not put routes in a single `routes.go`; each resource gets its own `*.route.go` and is wired up in `routes/router.go`.

## Backend architecture

- **Stack**: Go 1.26, Gin, GORM + PostgreSQL, JWT (golang-jwt/v5) + argon2id passwords, Cobra CLI, Swagger via swaggo, n8n integration.
- **Entry point**: `main.go` → `cmd.Execute()` (Cobra). Commands: `serve`, `migrate` (`--drop` to drop first).
- **Dependency injection**: `go.uber.org/dig`. Register every new repo/service/handler in `internal/container/di.go`. Repos and services are registered as interfaces via `dig.As(new(repo.IFooRepo))`. The `App` struct in `internal/app/app.go` receives only the top-level handlers; when adding a new domain you must: add it to `App`, update `app.New`, register in `di.go`, and wire its routes in `routes/router.go`.
- **Generic CRUD contract**: repos implement `common.IBaseCrudRepo[TEntity, TInput, TPage]` and services implement `common.IBaseCrudService[...]` (see `internal/common/base_crud.go`). Extend the interface for domain-specific methods (e.g. `Login`, `GetPermissions`).
- **GORM generics API**: this codebase uses the Go 1.26 `gorm.G[T](db)` form (e.g. `gorm.G[entities.User](r.db).Where(...).First(c)`), passing the `context.Context` (`c` / `ctx`) as the first arg. Don't switch to the non-generic `db.First(&model)` style.
- **Response envelope**: always return JSON via `common.ResultOk[T](items, msg)` / `common.ResultErr(err, msg)` (defined in `internal/common/result.go`), not raw `gin.H`.
- **Swagger annotations**: handlers carry swaggo `// @...` comments. Swagger docs live in `docs/` (gitignored, generated). Regeneration requires `make swag` — do not run it yourself; flag it to the user.
- **Config/env**: loaded via godotenv from `.env` in `internal/configs/env.go`. Required vars are in `example.env` (PORT, DB_*, JWT_SECRET, SERVER_ENV, DEFAULT_*, N8N_BASE_URL). `.env` is gitignored.

## Frontend (`internal/web/`)

- **Stack**: HTMX + Alpine.js (server-rendered HTML from Go, no SPA framework). Prefer HTMX partial swaps (`hx-get`/`hx-post`/`hx-swap`, OOB swaps, `hx-target`) over full page reloads or client-side state. Verify exact HTMX/Alpine attributes via Context7 before writing them.
- **Structure** (mirror the backend's per-domain file convention):
  - `internal/web/routes/` — `<name>.route.go`, wired up centrally (note: `internal/app/app.go` currently has the web route call commented out — re-enable when implementing the first real route).
  - `internal/web/handlers/` — `<name>.handler.go`.
  - `internal/web/viewmodels/` — `<name>.viewmodel.go`; put all related view models for a page in one file.
  - `internal/web/views/` — `.html` templates (`layout.html` is the base layout).
  - `internal/web/public/` — `css/` and `scripts/` static assets served directly.
- When adding a frontend page for an existing backend resource, reuse the backend services (don't reimplement data access in the web layer).

## Verification you can do

- `go build ./...` — confirm it compiles. Run this after edits; it's the safe check.
- `go vet ./...` — static checks.
- `go test ./...` — there is currently no test suite; don't assume one exists.
- Do **not** run `make dev`, `make go`, `air`, `serve`, or `migrate` to verify — see Hard rules.

## Other notes

- `.env`, `docs/`, and `tmp/` are gitignored. Don't commit `.env` or regenerated swagger.
- CI (`.github/workflows/docker-publish.yml`) builds and publishes the Docker image; deploy uses `docker-compose.yml` (runs `migrate` then `serve` in containers).
- Some code has commented-out role/permission logic (`entities.Role`, `UserRole`, `AssignRoles`) — it's intentionally stubbed, don't "fix" it unless asked.
