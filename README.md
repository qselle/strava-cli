# strava-cli

A command-line interface for the Strava API. Track your activities, monitor your streak, and let your AI agents keep you accountable.

## Install

```bash
go install github.com/qselle/strava-cli@latest
```

Or download a binary from [Releases](https://github.com/qselle/strava-cli/releases).

## Setup

1. Create a Strava API application at https://www.strava.com/settings/api
2. Authenticate (one-time):

```bash
STRAVA_CLIENT_ID=your_client_id STRAVA_CLIENT_SECRET=your_client_secret strava-cli auth
```

This prints a URL to open in your browser. Authorize the app, then paste the code back in the terminal. Credentials are stored locally — no env vars needed after this.

On a machine with a browser, add `--browser` for automatic callback.

## Usage

```bash
# List recent activities
strava-cli activities
strava-cli activities --after 2026-01-01 --limit 20

# Show all-time and year-to-date stats
strava-cli stats

# Check your streak
strava-cli streak
strava-cli streak --days 30

# JSON output (for AI agents and scripts)
strava-cli streak --json
strava-cli activities --json
```

## AI Agent Usage

AI agents can discover all commands via `--help`:

```bash
strava-cli --help
strava-cli streak --help
```

Use `--json` for structured output that agents can parse directly.

## MCP Server

Exposes three tools: `get_activities`, `get_stats`, and `get_streak`.

You must authenticate first (see [Setup](#setup)), then start the server:

```bash
strava-cli serve                # stdio transport
strava-cli serve --http :8080   # HTTP/SSE transport
```

### Claude Code / Claude Desktop

Add to your MCP config:

```json
{
  "mcpServers": {
    "strava": {
      "command": "strava-cli",
      "args": ["serve"]
    }
  }
}
```

## Development

```bash
make build    # build to bin/
make test     # run tests
make lint     # run linter
```

## License

MIT
