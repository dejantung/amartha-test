package repository

import "context"

//go:generate mockgen -destination=../mocks/mock_billing_cache.go -package=mocks billing-engine/internal/billing/repository BillingCacheProvider
type BillingCacheProvider interface {
	Set(ctx context.Context, key string, value interface{}) error
	Get(ctx context.Context, key string) (interface{}, error)
}
