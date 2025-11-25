# Repository Guidelines

## Project Structure & Module Organization
- Root entry point is `main.go`, wiring Cobra commands and Zap logging.
- CLI commands live in `cmd/` (`add`, `show`, `chkcfg`, `rmcfg`, `root`), each file handling a subcommand.
- Shared helpers for scanning repos and reading/writing the tracking file sit in `lib/`.
- Logging setup is isolated in `log/log.go`.
- User state (list of tracked repos) is stored in `~/.gogitlocalstats`; keep paths one per line.

## Build, Test, and Development Commands
- `go build .` — compile the CLI locally; target Go 1.18+ per `go.mod`.
- `go run . --help` — run without installing to verify flags and output.
- `go install .` — install `visualizedGit` into your GOPATH/bin for end-to-end usage.
- `go test ./...` — run package tests when added; should be clean before pushing.
- `gofmt -w .` — format all Go code; required before commits.

## Coding Style & Naming Conventions
- Follow Go defaults: tabs for indentation, camel-cased exported identifiers, `err` for errors.
- Prefer Zap logging (`zap.L()`) over `fmt`/`log` inside the app.
- Keep command definitions concise; flag names match existing ones (`-p` for path, `-e` for email).
- Avoid vendor-specific paths in repo scans (`vendor`, `node_modules` already skipped in `lib`).

## Testing Guidelines
- Standard Go testing: files end with `_test.go`, functions `TestXxx(t *testing.T)`.
- Use table-driven tests for `lib` helpers (e.g., `GetWeeksInLastSixMon`, file parsing).
- Mock filesystem interactions with temp dirs; avoid touching real `~/.gogitlocalstats` in tests.
- Aim for `go test ./...` to pass without needing external Git repos.

## Commit & Pull Request Guidelines
- Commit history uses short prefixes (e.g., `feat:`, `fix:`, `modify:`); keep messages imperative and scoped.
- Keep commits small and focused; include relevant CLI samples when behavior changes.
- PRs should describe intent, impacted commands, and manual checks (`go test`, `go run . --help` output).
- Link issues when available and add screenshots/terminal snippets for user-facing changes.

## Security & Configuration Tips
- The tool reads local repos and writes only to `~/.gogitlocalstats`; review changes before committing paths.
- Avoid logging sensitive repo paths or emails; Zap JSON logs go to stdout by default.
- Ensure new code respects existing skips for large directories and does not traverse network mounts unintentionally.
