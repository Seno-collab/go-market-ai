package cache

import (
	"github.com/redis/go-redis/v9"
)

func ConnectRedis(dsn string) (*redis.Client, error) {
	opt, err := redis.ParseURL(dsn)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opt)
	return client, nil
}


