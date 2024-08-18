package repository

import (
	"billing-engine/pkg/logger"
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
)

//go:generate mockgen -destination=../mocks/mock_billing_cache.go -package=mocks billing-engine/internal/billing/repository BillingCacheProvider
type BillingCacheProvider interface {
	Set(ctx context.Context, key string, value interface{}) error
	Get(ctx context.Context, key string) (interface{}, error)
}

type redisCache struct {
	client *redis.Client
	log    logger.Logger
}

func (r redisCache) Set(ctx context.Context, key string, value interface{}) error {
	r.log.WithField("key", key).
		WithField("value", value).Info("[Set] setting key to cache")

	valJson, _ := json.Marshal(value)
	err := r.client.Set(ctx, key, valJson, 0).Err()
	if err != nil {
		r.log.WithField("error", err).Error("[Set] failed to set key to cache")
		return err
	}

	r.log.WithField("key", key).Info("[Set] key set to cache")
	return nil
}

func (r redisCache) Get(ctx context.Context, key string) (interface{}, error) {
	r.log.WithField("key", key).Info("[Get] getting key from cache")

	val, err := r.client.Get(ctx, key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		r.log.WithField("error", err).Error("[Get] failed to get key from cache")
		return nil, err
	} else if err != nil && errors.Is(err, redis.Nil) {
		return nil, nil
	}

	r.log.WithField("key", key).Info("[Get] key retrieved from cache")
	return val, nil
}

func NewBillingCacheProvider(client *redis.Client, log logger.Logger) BillingCacheProvider {
	return &redisCache{
		client: client,
		log:    log,
	}
}
