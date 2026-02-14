package database

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

var RDB *redis.Client
var ctx = context.Background()

type RedisConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

func InitRedis(cfg *RedisConfig) error {
	RDB = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	_, err := RDB.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("连接 Redis 失败: %w", err)
	}
	return nil
}

func GetRedis() *redis.Client {
	return RDB
}

// Session 相关方法
const SessionPrefix = "session:"
const SessionExpire = 24 * time.Hour

func SetSession(token string, userID uint) error {
	return RDB.Set(ctx, SessionPrefix+token, userID, SessionExpire).Err()
}

func GetSession(token string) (uint, error) {
	val, err := RDB.Get(ctx, SessionPrefix+token).Uint64()
	if err != nil {
		return 0, err
	}
	return uint(val), nil
}

func DeleteSession(token string) error {
	return RDB.Del(ctx, SessionPrefix+token).Err()
}

func RefreshSession(token string) error {
	return RDB.Expire(ctx, SessionPrefix+token, SessionExpire).Err()
}
