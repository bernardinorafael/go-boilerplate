package rdb

import (
	"context"
	"fmt"

	"github.com/bernardinorafael/go-boilerplate/pkg/config"
	"github.com/bernardinorafael/go-boilerplate/pkg/fault"
	"github.com/redis/go-redis/v9"
)

type database struct {
	db *redis.Client
}

func NewConnection(ctx context.Context, config *config.Config) (*database, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.RedisHost, config.RedisPort),
		Password: config.RedisPassword,
		DB:       0,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fault.New("failed to connect to redis", fault.WithError(err))
	}

	return &database{db: client}, nil
}

func (r *database) Close() error {
	return r.db.Close()
}
