package cache

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"user-svc/helpers/fault"

	"github.com/redis/go-redis/v9"
)

// Set menyimpan value ke Redis dengan key dan TTL (time to live) tertentu.
// Jika gagal menyimpan, akan mengembalikan DetailedError bertipe ErrUnprocessable.
func Set(ctx context.Context, client *redis.Client, key string, value interface{}, ttl time.Duration) error {
	err := client.Set(ctx, key, value, ttl).Err()
	if err != nil {
		return fault.Custom(
			http.StatusUnprocessableEntity,
			fault.ErrUnprocessable,
			fmt.Sprintf("failed to save to Redis [key=%s]: %v", key, err),
		)
	}
	return nil
}

// Exist mengecek apakah key tertentu ada di Redis.
// Mengembalikan true jika key ditemukan, false jika tidak, atau error jika gagal mengecek.
func Exist(ctx context.Context, client *redis.Client, key string) (bool, error) {
	count, err := client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Get mengambil value dari Redis berdasarkan key.
// Mengembalikan value jika ditemukan, atau error jika gagal/get miss.
func Get(ctx context.Context, client *redis.Client, key string) (interface{}, error) {
	value, err := client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	return value, nil
}
