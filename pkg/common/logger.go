package common

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

type LoggerConfig struct {
	Level  string
	Pretty bool
	Output io.Writer
}

// Logger is the application-wide zerolog instance.
var Logger zerolog.Logger

func init() {
	Logger = zerolog.New(os.Stdout).With().Timestamp().Caller().Logger()
}

// NewLogger builds a zerolog.Logger using sensible defaults and stores it in Logger.
// - LOG_LEVEL environment variable controls the log level (defaults to info).
// - APP_ENV controls whether development-friendly console output is used.
func NewLogger(cfgs ...LoggerConfig) zerolog.Logger {
	cfg := buildConfig(cfgs...)

	zerolog.TimeFieldFormat = time.RFC3339Nano

	level := parseLogLevel(cfg.Level)
	logFile, err := os.OpenFile("application.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	logger := zerolog.New(logFile).
		Level(level).
		With().
		Timestamp().
		Logger()

	logger = logger.With().Caller().Logger()
	zerolog.SetGlobalLevel(level)

	Logger = logger
	return Logger
}

// FromContext extracts a zerolog logger stored in the context; falls back to the global Logger.
func FromContext(ctx context.Context) zerolog.Logger {
	if logger, ok := ctx.Value(loggerKey{}).(zerolog.Logger); ok {
		return logger
	}
	return Logger
}

// WithContext stores the provided logger inside the context for downstream handlers.
func WithContext(ctx context.Context, logger zerolog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}

type loggerKey struct{}

func buildConfig(cfgs ...LoggerConfig) LoggerConfig {
	defaultOutput := io.Writer(os.Stdout)
	pretty := strings.EqualFold(os.Getenv("APP_ENV"), "development")
	level := os.Getenv("LOG_LEVEL")

	cfg := LoggerConfig{
		Level:  level,
		Pretty: pretty,
		Output: defaultOutput,
	}
	if pretty {
		cfg.Output = zerolog.ConsoleWriter{
			Out:        defaultOutput,
			TimeFormat: time.RFC3339,
			FormatLevel: func(i interface{}) string {
				return fmt.Sprintf("[%-6s]", i) // Custom level format in square brackets
			},
			FormatMessage: func(i interface{}) string {
				return fmt.Sprintf("| %-20s |", i) // Custom message format surrounded by pipes
			},
			FormatCaller: func(i interface{}) string {
				return fmt.Sprintf("[%s]", i) // Custom caller format in brackets
			},
		}
	}

	if len(cfgs) > 0 {
		override := cfgs[0]
		if override.Level != "" {
			cfg.Level = override.Level
		}
		if override.Output != nil {
			cfg.Output = override.Output
		}
		if override.Pretty {
			cfg.Pretty = true
			cfg.Output = zerolog.ConsoleWriter{
				Out:        defaultOutput,
				TimeFormat: time.RFC3339,
				FormatLevel: func(i interface{}) string {
					return fmt.Sprintf("[%-6s]", i) // Custom level format in square brackets
				},
				FormatMessage: func(i interface{}) string {
					return fmt.Sprintf("| %-20s |", i) // Custom message format surrounded by pipes
				},
				FormatCaller: func(i interface{}) string {
					return fmt.Sprintf("[%s]", i) // Custom caller format in brackets
				},
			}
		}
	}

	return cfg
}

func parseLogLevel(level string) zerolog.Level {
	if level == "" {
		return zerolog.InfoLevel
	}
	parsedLevel, err := zerolog.ParseLevel(strings.ToLower(level))
	if err != nil {
		return zerolog.InfoLevel
	}
	return parsedLevel
}
