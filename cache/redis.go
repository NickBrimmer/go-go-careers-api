package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisCache(host, port string) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: "", // no password for local dev
		DB:       0,
	})

	ctx := context.Background()

	// Test connection
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisCache{
		client: client,
		ctx:    ctx,
	}, nil
}

// Get retrieves a value from cache and unmarshals it into the provided interface
func (c *RedisCache) Get(key string, dest interface{}) error {
	val, err := c.client.Get(c.ctx, key).Result()
	if err == redis.Nil {
		return fmt.Errorf("cache miss")
	} else if err != nil {
		return err
	}

	return json.Unmarshal([]byte(val), dest)
}

// Set stores a value in cache with a TTL
func (c *RedisCache) Set(key string, value interface{}, ttl time.Duration) error {
	json, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.client.Set(c.ctx, key, json, ttl).Err()
}

// Delete removes a key from cache
func (c *RedisCache) Delete(key string) error {
	return c.client.Del(c.ctx, key).Err()
}

// DeletePattern removes all keys matching a pattern
func (c *RedisCache) DeletePattern(pattern string) error {
	iter := c.client.Scan(c.ctx, 0, pattern, 0).Iterator()
	for iter.Next(c.ctx) {
		if err := c.client.Del(c.ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}
	return iter.Err()
}

// Close closes the Redis connection
func (c *RedisCache) Close() error {
	return c.client.Close()
}
