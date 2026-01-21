package logger

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/diode"
)

type LoggerConfig struct {
	Level              string
	Pretty             bool
	Output             io.Writer
	LogDir             string
	MaxSizeMB          int
	MaxBackupsPerDay   int
	RetentionDays      int
	AsyncQueue         int
	AsyncFlushInterval time.Duration
}

// Logger is the application-wide zerolog instance.
var Logger zerolog.Logger

const (
	defaultLogDir           = "logs"
	defaultLogRetentionDays = 7
	defaultLogMaxSizeMB     = 50
	defaultLogMaxBackups    = 5
	defaultAsyncQueue       = 1024
	defaultAsyncFlush       = 10 * time.Millisecond
)

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

	output := cfg.Output
	if output == nil {
		rotating, err := newRotatingWriter(rotationConfig{
			Dir:              cfg.LogDir,
			MaxSizeBytes:     int64(cfg.MaxSizeMB) * 1024 * 1024,
			RetentionDays:    cfg.RetentionDays,
			MaxBackupsPerDay: cfg.MaxBackupsPerDay,
		})
		if err != nil {
			panic(err)
		}
		output = rotating
	}

	if cfg.AsyncQueue > 0 {
		flushInterval := cfg.AsyncFlushInterval
		if flushInterval <= 0 {
			flushInterval = defaultAsyncFlush
		}
		output = diode.NewWriter(output, cfg.AsyncQueue, flushInterval, func(missed int) {
			fmt.Fprintf(os.Stderr, "logger queue full, dropped %d messages\n", missed)
		})
	}

	var writer io.Writer
	if cfg.Pretty {
		consoleWriter := zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
			FormatLevel: func(i any) string {
				return fmt.Sprintf("[%s]", i)
			},
			FormatMessage: func(i any) string {
				return fmt.Sprintf("| %s |", i)
			},
			FormatCaller: func(i any) string {
				return fmt.Sprintf("[%s]", i)
			},
		}
		writer = io.MultiWriter(consoleWriter, output)
	} else {
		writer = io.MultiWriter(os.Stdout, output)
	}

	logger := zerolog.New(writer).
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
	pretty := strings.EqualFold(os.Getenv("ENVIRONMENT"), "development")
	level := os.Getenv("LOG_LEVEL")
	cfg := LoggerConfig{
		Level:              level,
		Pretty:             pretty,
		LogDir:             envOrDefault("LOG_DIR", defaultLogDir),
		MaxSizeMB:          envInt("LOG_MAX_SIZE_MB", defaultLogMaxSizeMB),
		MaxBackupsPerDay:   envInt("LOG_MAX_BACKUPS_PER_DAY", defaultLogMaxBackups),
		RetentionDays:      envInt("LOG_RETENTION_DAYS", defaultLogRetentionDays),
		AsyncQueue:         envInt("LOG_ASYNC_QUEUE", defaultAsyncQueue),
		AsyncFlushInterval: defaultAsyncFlush,
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
		}
		if override.LogDir != "" {
			cfg.LogDir = override.LogDir
		}
		if override.MaxSizeMB > 0 {
			cfg.MaxSizeMB = override.MaxSizeMB
		}
		if override.MaxBackupsPerDay > 0 {
			cfg.MaxBackupsPerDay = override.MaxBackupsPerDay
		}
		if override.RetentionDays > 0 {
			cfg.RetentionDays = override.RetentionDays
		}
		if override.AsyncQueue >= 0 {
			cfg.AsyncQueue = override.AsyncQueue
		}
		if override.AsyncFlushInterval > 0 {
			cfg.AsyncFlushInterval = override.AsyncFlushInterval
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

type rotationConfig struct {
	Dir              string
	MaxSizeBytes     int64
	RetentionDays    int
	MaxBackupsPerDay int
}

type rotatingWriter struct {
	cfg         rotationConfig
	mu          sync.Mutex
	file        *os.File
	size        int64
	currentDate string
	sequence    int
}

func newRotatingWriter(cfg rotationConfig) (*rotatingWriter, error) {
	if cfg.Dir == "" {
		cfg.Dir = defaultLogDir
	}
	if cfg.RetentionDays <= 0 {
		cfg.RetentionDays = defaultLogRetentionDays
	}
	if cfg.MaxSizeBytes <= 0 {
		cfg.MaxSizeBytes = int64(defaultLogMaxSizeMB) * 1024 * 1024
	}
	if cfg.MaxBackupsPerDay <= 0 {
		cfg.MaxBackupsPerDay = defaultLogMaxBackups
	}

	if err := os.MkdirAll(cfg.Dir, 0755); err != nil {
		return nil, err
	}

	w := &rotatingWriter{cfg: cfg}
	if err := w.rotateIfNeeded(0); err != nil {
		return nil, err
	}
	return w, nil
}

func (w *rotatingWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if err := w.rotateIfNeeded(len(p)); err != nil {
		return 0, err
	}
	n, err := w.file.Write(p)
	if err != nil {
		return n, err
	}
	w.size += int64(n)
	return n, nil
}

func (w *rotatingWriter) rotateIfNeeded(nextWriteLen int) error {
	now := time.Now()
	today := now.Format("2006-01-02")

	if w.file == nil || today != w.currentDate {
		if err := w.openNewFile(today, 0); err != nil {
			return err
		}
		if err := w.cleanupRetention(); err != nil {
			fmt.Fprintf(os.Stderr, "log retention cleanup failed: %v\n", err)
		}
		return nil
	}

	if w.cfg.MaxSizeBytes > 0 && w.size+int64(nextWriteLen) > w.cfg.MaxSizeBytes {
		_ = w.file.Close()
		w.sequence++
		if err := w.ensureDailyLimit(today); err != nil {
			fmt.Fprintf(os.Stderr, "log backup cleanup failed: %v\n", err)
		}
		return w.openNewFile(today, w.sequence)
	}

	return nil
}

func (w *rotatingWriter) openNewFile(date string, seq int) error {
	if err := os.MkdirAll(w.cfg.Dir, 0755); err != nil {
		return err
	}

	path := logFilePath(w.cfg.Dir, date, seq)
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}

	if w.file != nil {
		_ = w.file.Close()
	}
	w.file = file
	w.currentDate = date
	w.sequence = seq

	info, err := file.Stat()
	if err != nil {
		w.size = 0
		return nil
	}
	w.size = info.Size()
	return nil
}

func (w *rotatingWriter) ensureDailyLimit(date string) error {
	if w.cfg.MaxBackupsPerDay <= 0 {
		return nil
	}

	files, err := os.ReadDir(w.cfg.Dir)
	if err != nil {
		return err
	}

	var matches []os.DirEntry
	for _, entry := range files {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".log" {
			continue
		}
		if strings.HasPrefix(entry.Name(), date) {
			matches = append(matches, entry)
		}
	}

	if len(matches)+1 <= w.cfg.MaxBackupsPerDay {
		return nil
	}

	sort.Slice(matches, func(i, j int) bool {
		ii, _ := matches[i].Info()
		ij, _ := matches[j].Info()
		if ii == nil || ij == nil {
			return matches[i].Name() < matches[j].Name()
		}
		return ii.ModTime().Before(ij.ModTime())
	})

	toRemove := len(matches) + 1 - w.cfg.MaxBackupsPerDay
	for i := 0; i < toRemove && i < len(matches); i++ {
		_ = os.Remove(filepath.Join(w.cfg.Dir, matches[i].Name()))
	}
	return nil
}

func (w *rotatingWriter) cleanupRetention() error {
	if w.cfg.RetentionDays <= 0 {
		return nil
	}

	cutoff := time.Now().AddDate(0, 0, -w.cfg.RetentionDays)
	entries, err := os.ReadDir(w.cfg.Dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".log" {
			continue
		}

		logDate, ok := parseLogDate(entry.Name())
		if !ok {
			continue
		}
		if logDate.Before(cutoff) {
			_ = os.Remove(filepath.Join(w.cfg.Dir, entry.Name()))
		}
	}
	return nil
}

func parseLogDate(name string) (time.Time, bool) {
	base := strings.TrimSuffix(name, filepath.Ext(name))
	parts := strings.Split(base, "-")
	if len(parts) < 3 {
		return time.Time{}, false
	}
	dateStr := strings.Join(parts[:3], "-")
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Time{}, false
	}
	return t, true
}

func logFilePath(dir, date string, seq int) string {
	if seq <= 0 {
		return filepath.Join(dir, fmt.Sprintf("%s.log", date))
	}
	return filepath.Join(dir, fmt.Sprintf("%s-%02d.log", date, seq))
}

func envOrDefault(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func envInt(key string, fallback int) int {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
	}
	return fallback
}
