package repository

import (
	"Dzaakk/simple-commerce/internal/auth/model"
	"context"
	"encoding/json"
	"os"

	"github.com/go-redis/redis/v8"
)

type AuthCacheSellerImpl struct {
	Client *redis.Client
}

func NewAuthCacheSellerRepository(client *redis.Client) AuthCacheSeller {
	return &AuthCacheSellerImpl{Client: client}
}

var (
	PrefixSellerToken        = os.Getenv("REDIS_PREFIX_SELLER")
	PrefixSellerRegistration = os.Getenv("REDIS_PREFIX_REGISTRATION_SELLER")
)

func (cache *AuthCacheSellerImpl) SetActivationSeller(c context.Context, email string, activationCode string) error {
	key := email + PrefixCodeActivation

	jsonData, err := json.Marshal(activationCode)
	if err != nil {
		return err
	}

	return cache.Client.Set(c, key, jsonData, ActivationExpired).Err()
}

func (cache *AuthCacheSellerImpl) GetActivationSeller(c context.Context, email string) (string, error) {
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

func (cache *AuthCacheSellerImpl) SetTokenSeller(c context.Context, email, token string) error {
	key := email + PrefixSellerToken

	jsonToken, err := json.Marshal(token)
	if err != nil {
		return err
	}

	return cache.Client.Set(c, key, jsonToken, TokenExpired).Err()
}

func (cache *AuthCacheSellerImpl) GetTokenSeller(c context.Context, email string) (*string, error) {
	key := email + PrefixSellerToken

	val, err := cache.Client.Get(c, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var token string
	err = json.Unmarshal([]byte(val), &token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (cache *AuthCacheSellerImpl) SetSellerRegistration(c context.Context, data model.SellerRegistrationReq) error {
	key := data.Email + PrefixSellerRegistration

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return cache.Client.Set(c, key, jsonData, ActivationExpired).Err()
}

func (cache *AuthCacheSellerImpl) GetSellerRegistration(c context.Context, email string) (*model.SellerRegistrationReq, error) {
	key := email + PrefixSellerRegistration

	val, err := cache.Client.Get(c, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var data model.SellerRegistrationReq
	err = json.Unmarshal([]byte(val), &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (cache *AuthCacheSellerImpl) DeleteTokenSeller(c context.Context, email string) error {
	key := email + PrefixSellerToken

	return cache.Client.Del(c, key).Err()
}
