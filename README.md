# Clipboard Manager

A lightweight, privacy-focused clipboard manager built in Go with clean architecture principles.

## Features

- **Background Monitoring** - Automatically tracks clipboard changes
- **Persistent Storage** - SQLite database stores complete clipboard history
- **Privacy Protection** - Detects and skips sensitive data (passwords, tokens, API keys)
- **HTTP API** - RESTful API over Unix socket for secure, local-only access
- **CLI Tool** - Command-line interface to query, search, and manage history
- **Structured Logging** - Configurable JSON/text logs with automatic rotation
- **Graceful Shutdown** - Proper cleanup and resource management
- **Comprehensive Tests** - Well-tested service layer with table-driven tests

## Quick Start

```bash
# Build
make build

# Start daemon
./bin/clipd

# Use CLI
./bin/clipctl list           # View recent entries
./bin/clipctl search "text"  # Search history
./bin/clipctl stats          # Show statistics
```

## Architecture

Built with clean architecture principles:

- **Layered Design** - Clear separation between API, business logic, and infrastructure
- **Dependency Injection** - Components communicate through interfaces
- **Unix Philosophy** - Does one thing well, uses Unix sockets for security
- **Testable** - Service layer tested independently of infrastructure

**Note:** This project is also a learning exercise in Go and software design patterns.

## Configuration

Configure via command-line flags:

```bash
./bin/clipd \
  --db ./clipboard.db \
  --socket /tmp/clipd.sock \
  --poll-interval 500ms \
  --log-level info \
  --log-format text
```

### Available Flags

| Flag              | Description                       | Default            |
| ----------------- | --------------------------------- | ------------------ |
| `--db`            | Database path                     | `./clipboard.db`   |
| `--socket`        | Unix socket path                  | `/tmp/clipd.sock`  |
| `--poll-interval` | Clipboard check interval          | `500ms`            |
| `--log-level`     | Log level (debug/info/warn/error) | `info`             |
| `--log-format`    | Log format (text/json)            | `text`             |
| `--log-output`    | Log output (stdout/file/both)     | `both`             |
| `--log-file`      | Log file path                     | `./logs/clipd.log` |

## API Reference

### Endpoints

```bash
# Health check
curl --unix-socket /tmp/clipd.sock http://unix/api/health

# List history
curl --unix-socket /tmp/clipd.sock http://unix/api/v1/history?limit=10

# Get specific entry
curl --unix-socket /tmp/clipd.sock http://unix/api/v1/history/1

# Search
curl --unix-socket /tmp/clipd.sock "http://unix/api/v1/search?q=example"

# Delete entry
curl --unix-socket /tmp/clipd.sock -X DELETE http://unix/api/v1/history/1

# Statistics
curl --unix-socket /tmp/clipd.sock http://unix/api/v1/stats
```

## Security & Privacy

- **Sensitive Data Detection** - Automatically filters passwords, tokens, and API keys
- **Unix Socket** - API only accessible locally (not over network)
- **File Permissions** - Socket has 0600 permissions (owner-only)
- **No Cloud** - Everything stays on your machine

## Project Structure

```
clipboard-manager/
├── cmd/
│   ├── clipd/           # Daemon
│   └── clipctl/         # CLI tool
├── internal/
│   ├── api/             # HTTP handlers and routes
│   ├── service/         # Business logic
│   ├── storage/         # Database layer
│   ├── daemon/          # Orchestration
│   ├── monitor/         # Clipboard monitoring
│   ├── analyzer/        # Content analysis
│   ├── config/          # Configuration
│   └── logger/          # Logging
└── Makefile
```

## Roadmap

**Current Status:**

- ✅ Core clipboard management
- ✅ API and CLI interface
- ✅ Structured logging
- ✅ Configuration system

**Planned Features:**

- 🔄 Automatic retention policy (cleanup old entries)
- 🔄 Error resilience (retry logic, circuit breakers)
- 📋 TUI (Terminal User Interface)
- 📋 Entry pinning and favorites
- 📋 Systemd service integration

## License

MIT License - See LICENSE file for details
