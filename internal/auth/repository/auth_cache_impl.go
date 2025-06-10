package repository

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

type AuthCache struct {
	Client *redis.Client
}

func NewAuthCacheRepository(client *redis.Client) AuthCacheRepository {
	return &AuthCache{Client: client}
}

const (
	CodeActivationExpired = time.Minute * 15
)

var (
	PrefixCodeActivation = os.Getenv("REDIS_PREFIX_CODE")
)

func (cache *AuthCache) SetActivationCustomer(c context.Context, email string, activationCode string) error {
	key := email + PrefixCodeActivation

	jsonData, err := json.Marshal(activationCode)
	if err != nil {
		return err
	}

	return cache.Client.Set(c, key, jsonData, CodeActivationExpired).Err()
}

func (cache *AuthCache) GetActivationCustomer(c context.Context, email string) (string, error) {
	key := email + PrefixCodeActivation

	val, err := cache.Client.Get(c, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
		return "", err
	}

	var activationCode string
	err = json.Unmarshal([]byte(val), &activationCode)
	if err != nil {
		return "", err
	}

	return activationCode, nil
}
