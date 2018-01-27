package Core

import (
	"crypto/md5"
	"encoding/hex"

	"strings"

	"fmt"

	"github.com/go-redis/redis"
)

type redisWrapper struct {
	client *redis.Client
}

type RedisKeyExistsError string

func (key RedisKeyExistsError) Error() string {
	return fmt.Sprintf("Key %v Exists", string(key))
}

func (redisWrapper) GenerateKeyHash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))

}
func (wrapper redisWrapper) InsertValue(value interface{}, keyHashParameters ...string) (string, error) {
	hash := wrapper.GenerateKeyHash(strings.Join(keyHashParameters, ""))
	key, err := wrapper.client.Exists(hash).Result()
	if err != nil {
		return "", err
	}
	if key == 0 {
		wrapper.client.Set(hash, value, 0)
		return hash, nil
	}

	return "", RedisKeyExistsError(hash)

}
func (wrapper redisWrapper) GetValue(key string) (interface{}, error) {
	return wrapper.client.Get(key).Bytes()
}

func NewRedisWrapper() redisWrapper {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0})
	return redisWrapper{client: client}
}
