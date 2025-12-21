package container

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	postgresMaxRetries     = 8
	postgresInitialBackoff = time.Second
	postgresDialTimeout    = 15 * time.Second
)

func ConnectPostgres(dsn string) (*pgxpool.Pool, error) {
	var pool *pgxpool.Pool
	var err error

	backoff := postgresInitialBackoff
	for attempt := 1; attempt <= postgresMaxRetries; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), postgresDialTimeout)

		pool, err = pgxpool.New(ctx, dsn)
		if err == nil {
			if pingErr := pool.Ping(ctx); pingErr == nil {
				cancel()
				return pool, nil
			} else {
				err = pingErr
			}
			pool.Close()
		}
		cancel()

		if attempt == postgresMaxRetries {
			break
		}

		time.Sleep(backoff)
		backoff *= 2
	}

	return nil, fmt.Errorf("cannot connect to postgres after %d attempts: %w", postgresMaxRetries, err)
}
