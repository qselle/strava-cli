# strava-cli

Go CLI for the Strava API with optional MCP server support.

## Build & Run

```bash
make build          # builds to bin/strava-cli
make test           # runs all tests
make lint           # runs golangci-lint
make install        # go install to $GOPATH/bin
```

Requires Go 1.22+.

## Architecture

- `main.go` — entry point, calls `cmd.Execute()`
- `cmd/` — Cobra command definitions (thin wiring, calls into `internal/`)
- `internal/auth/` — OAuth2 flow, token storage and refresh
- `internal/api/` — Strava API HTTP client
- `internal/server/` — MCP server wrapper (not yet implemented)

## Environment Variables

- `STRAVA_CLIENT_ID` — required, from https://www.strava.com/settings/api
- `STRAVA_CLIENT_SECRET` — required, same source

## Tokens

Stored at `~/.config/strava-cli/token.json`. Access tokens auto-refresh using the refresh token.

## Conventions

- Follow standard Go project layout (`cmd/` + `internal/`)
- Use Cobra for all CLI commands
- All output commands support `--json` flag for structured output
- Keep commands thin — business logic belongs in `internal/`
