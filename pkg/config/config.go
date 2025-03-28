package config

import (
	"flag"
	"os"
	"strconv"
	"strings"

	"github.comm/kevinlucasklein/madden-discord-bot/pkg/utils"
)

// Config holds application configuration
type Config struct {
	Port      int
	ExportURL string
	DataDir   string
	LogLevel  utils.LogLevel
	LogToFile bool
	LogDir    string
}

// Default configuration values
const (
	DefaultPort      = 8080
	DefaultExportURL = "/export"
	DefaultDataDir   = "./data"
	DefaultLogLevel  = utils.LogLevelDebug
	DefaultLogToFile = true
	DefaultLogDir    = "./logs"
)

// Load loads configuration from environment variables and command-line flags
func Load() *Config {
	config := &Config{
		Port:      DefaultPort,
		ExportURL: DefaultExportURL,
		DataDir:   DefaultDataDir,
		LogLevel:  DefaultLogLevel,
		LogToFile: DefaultLogToFile,
		LogDir:    DefaultLogDir,
	}

	// Load from environment variables first
	if portStr := os.Getenv("MADDEN_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			config.Port = port
		}
	}
	if exportURL := os.Getenv("MADDEN_EXPORT_URL"); exportURL != "" {
		config.ExportURL = exportURL
	}
	if dataDir := os.Getenv("MADDEN_DATA_DIR"); dataDir != "" {
		config.DataDir = dataDir
	}
	if logLevel := os.Getenv("MADDEN_LOG_LEVEL"); logLevel != "" {
		config.LogLevel = parseLogLevel(logLevel)
	}
	if logToFile := os.Getenv("MADDEN_LOG_TO_FILE"); logToFile != "" {
		config.LogToFile = strings.ToLower(logToFile) == "true"
	}
	if logDir := os.Getenv("MADDEN_LOG_DIR"); logDir != "" {
		config.LogDir = logDir
	}

	// Command-line flags override environment variables
	port := flag.Int("port", config.Port, "Port for the HTTP server")
	exportURL := flag.String("export-url", config.ExportURL, "URL path for receiving exports")
	dataDir := flag.String("data-dir", config.DataDir, "Directory to store export data")
	logLevelStr := flag.String("log-level", logLevelToString(config.LogLevel), "Log level (debug, info, warn, error)")
	logToFile := flag.Bool("log-to-file", config.LogToFile, "Whether to log to a file")
	logDir := flag.String("log-dir", config.LogDir, "Directory to store log files")
	flag.Parse()

	// Override with command-line values if specified
	config.Port = *port
	config.ExportURL = *exportURL
	config.DataDir = *dataDir
	config.LogLevel = parseLogLevel(*logLevelStr)
	config.LogToFile = *logToFile
	config.LogDir = *logDir

	return config
}

// parseLogLevel converts a string log level to LogLevel
func parseLogLevel(level string) utils.LogLevel {
	switch strings.ToLower(level) {
	case "debug":
		return utils.LogLevelDebug
	case "info":
		return utils.LogLevelInfo
	case "warn", "warning":
		return utils.LogLevelWarn
	case "error":
		return utils.LogLevelError
	default:
		return DefaultLogLevel
	}
}

// logLevelToString converts a LogLevel to string
func logLevelToString(level utils.LogLevel) string {
	switch level {
	case utils.LogLevelDebug:
		return "debug"
	case utils.LogLevelInfo:
		return "info"
	case utils.LogLevelWarn:
		return "warn"
	case utils.LogLevelError:
		return "error"
	default:
		return "info"
	}
}
