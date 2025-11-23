# Repository Guidelines

## Project Structure & Module Organization
- CLI entrypoint and commands live in `cmd/clash-fish/` (e.g., `main.go`, `start.go`, `config.go`); add new commands here.
- Shared utilities sit in `pkg/` (`logger`, `constants`, `utils`). Keep helpers generic and dependency-light.
- Core domain packages are under `internal/` (e.g., `config`, `proxy`, `system`, `subscription`, `service`); prefer adding new behavior here rather than `cmd/`.
- Config assets belong in `configs/` (rules, templates). Scripts and automation go in `scripts/`. Build artifacts output to `build/`.

## Build, Test, and Development Commands
- `make build` — compile the CLI to `build/clash-fish`.
- `make run` — build then run `clash-fish start` with sudo (for TUN/root needs).
- `make test` — run `go test -v ./...` across all modules.
- `make fmt` — run `go fmt ./...` to enforce formatting.
- `make install` / `make uninstall` — install or remove the binary in `/usr/local/bin`.
- Quick manual build: `go build -o build/clash-fish ./cmd/clash-fish`.

## Coding Style & Naming Conventions
- Go 1.21+; always `go fmt` before pushing. Prefer standard library first, then third-party imports, grouped and sorted.
- Keep files ASCII; log and user-facing text in English unless a feature mandates otherwise.
- Use lowercase-hyphen command names (`clash-fish`) and lowerCamelCase for Go locals; exported identifiers should be PascalCase with concise doc comments where behavior isn’t obvious.
- Avoid magic strings/paths; centralize defaults in `pkg/constants`.

## Testing Guidelines
- Add unit tests alongside implementation packages under `internal/<module>/`; name files `*_test.go`.
- Prefer table-driven tests for command handlers and utilities. Use fakes over network calls; keep tests hermetic (no external services).
- Run `make test` before opening a PR; add coverage for edge cases (permission errors, missing config, already-running service).

## Commit & Pull Request Guidelines
- Commits: concise, imperative subject (e.g., `Add start command stub`, `Fix log path creation`). Squash noisy WIP commits before review when possible.
- PRs: include what changed, why, and how to verify (commands run, screenshots not required for CLI). Link issues if available; call out risk areas (permissions, file paths) and follow-up TODOs.

## Security & Configuration Tips
- Avoid checking in user configs or logs; `.gitignore` already excludes `config.yaml`, `profiles/`, and `logs/`.
- Actions requiring root should validate with `pkg/utils.CheckRoot`; avoid silent elevation.
- Default config/log dirs resolve to `~/.config/clash-fish/`; keep new assets under that root unless a flag overrides it.
