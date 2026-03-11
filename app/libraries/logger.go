package libraries

import (
	"io"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

var Logger zerolog.Logger

// LogConfig holds logger configuration
type LogConfig struct {
	Level  string // debug, info, warn, error
	Format string // json, console
}

// InitLogger initializes logger with default settings (console, info level)
func InitLogger() zerolog.Logger {
	return InitLoggerWithConfig(LogConfig{
		Level:  getEnvOrDefault("LOG_LEVEL", "info"),
		Format: getEnvOrDefault("LOG_FORMAT", "console"),
	})
}

// InitLoggerWithConfig initializes logger with specific configuration
func InitLoggerWithConfig(cfg LogConfig) zerolog.Logger {
	// Set output format
	var output io.Writer
	if strings.ToLower(cfg.Format) == "json" {
		output = os.Stdout
	} else {
		output = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	}

	// Set log level
	level := parseLogLevel(cfg.Level)

	logger := zerolog.New(output).
		Level(level).
		With().
		Timestamp().
		Caller().
		Logger()

	Logger = logger
	return logger
}

func parseLogLevel(level string) zerolog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn", "warning":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	default:
		return zerolog.InfoLevel
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func GetLogger() zerolog.Logger {
	if Logger.GetLevel() == zerolog.NoLevel {
		InitLogger()
	}
	return Logger
}

// With creates a sublogger with additional context
func With() zerolog.Context {
	return Logger.With()
}
