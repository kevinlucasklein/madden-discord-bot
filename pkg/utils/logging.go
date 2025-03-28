package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// LogLevel represents logging level
type LogLevel int

const (
	// LogLevelDebug is for debug messages
	LogLevelDebug LogLevel = iota
	// LogLevelInfo is for informational messages
	LogLevelInfo
	// LogLevelWarn is for warning messages
	LogLevelWarn
	// LogLevelError is for error messages
	LogLevelError
)

// Logger provides logging functionality
type Logger struct {
	level      LogLevel
	debugLog   *log.Logger
	infoLog    *log.Logger
	warnLog    *log.Logger
	errorLog   *log.Logger
	logFile    *os.File
	enableFile bool
}

// NewLogger creates a new logger instance
func NewLogger(level LogLevel, logToFile bool, logDir string) (*Logger, error) {
	// Create standard loggers
	debugLog := log.New(os.Stdout, "[DEBUG] ", log.LstdFlags)
	infoLog := log.New(os.Stdout, "[INFO]  ", log.LstdFlags)
	warnLog := log.New(os.Stdout, "[WARN]  ", log.LstdFlags)
	errorLog := log.New(os.Stderr, "[ERROR] ", log.LstdFlags)

	logger := &Logger{
		level:      level,
		debugLog:   debugLog,
		infoLog:    infoLog,
		warnLog:    warnLog,
		errorLog:   errorLog,
		enableFile: logToFile,
	}

	// Set up file logging if enabled
	if logToFile {
		if err := EnsureDirectoryExists(logDir); err != nil {
			return nil, fmt.Errorf("failed to create log directory: %w", err)
		}

		// Create log file with timestamp
		timestamp := time.Now().Format("20060102")
		logFilePath := filepath.Join(logDir, fmt.Sprintf("madden_%s.log", timestamp))
		logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}

		// Update loggers to write to both console and file
		logger.logFile = logFile
		logger.debugLog = log.New(os.Stdout, "[DEBUG] ", log.LstdFlags|log.Lshortfile)
		logger.infoLog = log.New(os.Stdout, "[INFO]  ", log.LstdFlags)
		logger.warnLog = log.New(os.Stdout, "[WARN]  ", log.LstdFlags)
		logger.errorLog = log.New(os.Stdout, "[ERROR] ", log.LstdFlags|log.Lshortfile)
	}

	return logger, nil
}

// Close closes the log file if it was opened
func (l *Logger) Close() {
	if l.logFile != nil {
		l.logFile.Close()
	}
}

// Debug logs a debug message
func (l *Logger) Debug(format string, v ...interface{}) {
	if l.level <= LogLevelDebug {
		l.debugLog.Printf(format, v...)
		if l.enableFile && l.logFile != nil {
			fmt.Fprintf(l.logFile, "[DEBUG] "+format+"\n", v...)
		}
	}
}

// Info logs an informational message
func (l *Logger) Info(format string, v ...interface{}) {
	if l.level <= LogLevelInfo {
		l.infoLog.Printf(format, v...)
		if l.enableFile && l.logFile != nil {
			fmt.Fprintf(l.logFile, "[INFO]  "+format+"\n", v...)
		}
	}
}

// Warn logs a warning message
func (l *Logger) Warn(format string, v ...interface{}) {
	if l.level <= LogLevelWarn {
		l.warnLog.Printf(format, v...)
		if l.enableFile && l.logFile != nil {
			fmt.Fprintf(l.logFile, "[WARN]  "+format+"\n", v...)
		}
	}
}

// Error logs an error message
func (l *Logger) Error(format string, v ...interface{}) {
	if l.level <= LogLevelError {
		l.errorLog.Printf(format, v...)
		if l.enableFile && l.logFile != nil {
			fmt.Fprintf(l.logFile, "[ERROR] "+format+"\n", v...)
		}
	}
}
