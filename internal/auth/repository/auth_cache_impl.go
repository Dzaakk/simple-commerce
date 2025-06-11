package repository

import (
	"Dzaakk/simple-commerce/internal/auth/model"
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
	ActivationExpired = time.Minute * 15
	TokenExpired      = time.Minute * 30
)

var (
	PrefixCodeActivation       = os.Getenv("REDIS_PREFIX_CODE")
	PrefixCustomerToken        = os.Getenv("REDIS_PREFIX_CUSTOMER")
	PrefixCustomerRegistration = os.Getenv("REDIS_PREFIX_REGISTRATION_CUSTOMER")
)

func (cache *AuthCache) SetActivationCustomer(c context.Context, email string, activationCode string) error {
	key := email + PrefixCodeActivation

	jsonData, err := json.Marshal(activationCode)
	if err != nil {
		return err
	}

	return cache.Client.Set(c, key, jsonData, ActivationExpired).Err()
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

func (cache *AuthCache) SetTokenCustomer(c context.Context, data model.CustomerToken) error {
	key := data.Email + PrefixCustomerToken

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return cache.Client.Set(c, key, jsonData, TokenExpired).Err()
}

func (cache *AuthCache) GetTokenCustomer(c context.Context, email string) (*model.CustomerToken, error) {
	key := email + PrefixCustomerToken

	val, err := cache.Client.Get(c, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var token model.CustomerToken
	err = json.Unmarshal([]byte(val), &token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (cache *AuthCache) SetCustomerRegistration(c context.Context, data model.CustomerRegistrationReq) error {
	key := data.Email + PrefixCustomerRegistration

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return cache.Client.Set(c, key, jsonData, ActivationExpired).Err()
}

func (cache *AuthCache) GetCustomerRegistration(c context.Context, email string) (*model.CustomerRegistrationReq, error) {
	key := email + PrefixCustomerRegistration

	val, err := cache.Client.Get(c, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var data model.CustomerRegistrationReq
	err = json.Unmarshal([]byte(val), &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
