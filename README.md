# Madden Discord Bot

A service to handle data exports from the Madden Companion App and integrate with Discord.

## Features

- HTTP server to receive exports from the Madden Companion App
- Stores export data as JSON files
- Configurable via environment variables or command-line flags

## Requirements

- Go 1.16 or higher

## Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/madden-discord-bot.git
cd madden-discord-bot

# Build the application
go build -o madden-bot
```

## Usage

```bash
# Run with default settings
./madden-bot

# Run with custom port
./madden-bot -port 9090

# Run with custom export URL and data directory
./madden-bot -export-url /madden-export -data-dir ./madden-data
```

### Environment Variables

You can also configure the application using environment variables:

- `MADDEN_PORT`: HTTP server port (default: 8080)
- `MADDEN_EXPORT_URL`: Export endpoint URL path (default: /export)
- `MADDEN_DATA_DIR`: Directory to store export data (default: ./data)

### Madden Companion App Setup

1. Open the Madden Companion App on your mobile device
2. Go to "Export"
3. Enter your server URL: `http://your-server-ip:8080/export`
4. Select the league and data you want to export
5. Press the export button

## Project Structure

```
madden-discord-bot/
├── main.go              # Application entry point
├── pkg/
│   ├── config/          # Configuration handling
│   │   └── config.go
│   └── madden/          # Madden service implementation
│       ├── handlers.go  # HTTP handlers
│       ├── models.go    # Data models
│       └── service.go   # Core service logic
├── data/                # Default directory for exported data
└── README.md            # This file
```

## Future Plans

- Discord integration for notifying users about new exports
- Data processing and statistics generation
- Web interface for viewing export data

## License

MIT
