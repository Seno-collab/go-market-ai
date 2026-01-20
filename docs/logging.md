# Production Logging

Features:
- JSON logs with `time`, `level`, `message`, `request_id` (from middleware), caller.
- Rotation by day and by size (per-day suffix `-01`, `-02`, ...).
- Retention policy to delete old log files automatically.
- Async, non-blocking writes using zerolog diode (drops with warning if queue is full).

## Defaults
- Directory: `logs/`
- Rotation: daily + max size `50MB` (`LOG_MAX_SIZE_MB`)
- Max files per day: `5` (`LOG_MAX_BACKUPS_PER_DAY`) to avoid file explosion
- Retention: `7` days (`LOG_RETENTION_DAYS`)
- Async queue: `1024` messages (`LOG_ASYNC_QUEUE`), flush every `10ms`
- Level: `info` (`LOG_LEVEL`)

## Environment variables
```
LOG_LEVEL=info
LOG_DIR=logs
LOG_MAX_SIZE_MB=50
LOG_MAX_BACKUPS_PER_DAY=5
LOG_RETENTION_DAYS=7
LOG_ASYNC_QUEUE=1024      # set 0 to disable async
ENVIRONMENT=production    # use "development" to enable console pretty logs
```

## Usage
```go
import "go-ai/pkg/logger"

func main() {
    log := logger.NewLogger()
    log.Info().Str("request_id", "abc-123").Msg("service started")
}
```

Override in code if you prefer:
```go
log := logger.NewLogger(logger.LoggerConfig{
    Level:            "debug",
    LogDir:           "/var/log/myapp",
    MaxSizeMB:        100,
    MaxBackupsPerDay: 10,
    RetentionDays:    14,
    AsyncQueue:       2048,
})
```

## Behavior notes
- Files are named `YYYY-MM-DD.log`, then `YYYY-MM-DD-01.log`, `-02`... when the size cap is exceeded in the same day.
- When the daily limit is reached, oldest files for that date are deleted before opening a new one.
- Retention uses the date in the filename; files older than `LOG_RETENTION_DAYS` are removed on the first write of a new day.
- Async buffer drops oldest messages when full and prints a warning to stderr. Increase `LOG_ASYNC_QUEUE` if you see drops.
