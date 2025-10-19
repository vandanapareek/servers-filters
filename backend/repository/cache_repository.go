package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisCacheRepository implements CacheRepository for Redis
type RedisCacheRepository struct {
	client *redis.Client
}

// Create new Redis cache repository
func NewRedisCacheRepository(client *redis.Client) CacheRepository {
	return &RedisCacheRepository{client: client}
}

// Get value from cache
func (r *RedisCacheRepository) Get(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil // Key not found, not an error
		}
		return "", fmt.Errorf("failed to get from cache: %w", err)
	}
	return val, nil
}

// Set value in cache with TTL
func (r *RedisCacheRepository) Set(ctx context.Context, key string, value string, ttl int) error {
	err := r.client.Set(ctx, key, value, time.Duration(ttl)*time.Second).Err()
	if err != nil {
		return fmt.Errorf("failed to set cache: %w", err)
	}
	return nil
}

// Delete value from cache
func (r *RedisCacheRepository) Delete(ctx context.Context, key string) error {
	err := r.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete from cache: %w", err)
	}
	return nil
}

// Clear all values from cache
func (r *RedisCacheRepository) Clear(ctx context.Context) error {
	err := r.client.FlushDB(ctx).Err()
	if err != nil {
		return fmt.Errorf("failed to clear cache: %w", err)
	}
	return nil
}

// Return Unmarshal JSON value from cache
func (r *RedisCacheRepository) GetJSON(ctx context.Context, key string, dest interface{}) error {
	val, err := r.Get(ctx, key)
	if err != nil {
		return err
	}
	if val == "" {
		return nil // Key not found
	}

	err = json.Unmarshal([]byte(val), dest)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON from cache: %w", err)
	}
	return nil
}

// Marshals and stores a JSON value in cache with TTL
func (r *RedisCacheRepository) SetJSON(ctx context.Context, key string, value interface{}, ttl int) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON for cache: %w", err)
	}

	return r.Set(ctx, key, string(data), ttl)
}

// NoOpCacheRepository implements CacheRepository as a no-op (disabled cache)
type NoOpCacheRepository struct{}

// NewNoOpCacheRepository creates a new no-op cache repository
func NewNoOpCacheRepository() CacheRepository {
	return &NoOpCacheRepository{}
}

// Get empty string (cache miss)
func (r *NoOpCacheRepository) Get(ctx context.Context, key string) (string, error) {
	return "", nil
}

// Set does nothing
func (r *NoOpCacheRepository) Set(ctx context.Context, key string, value string, ttl int) error {
	return nil
}

// Delete does nothing
func (r *NoOpCacheRepository) Delete(ctx context.Context, key string) error {
	return nil
}

// Clear does nothing
func (r *NoOpCacheRepository) Clear(ctx context.Context) error {
	return nil
}
