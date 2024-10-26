package inventory

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
	model "inventory-service/internal/model/inventory"
)

type RedisRepository interface {
	CacheProduct(ctx context.Context, product *model.Product) error
	GetCachedProduct(ctx context.Context, id string) (*model.Product, error)
	DeleteCachedProduct(ctx context.Context, id string) error
}

type redisRepository struct {
	redisClient *redis.Client
}

var _ RedisRepository = (*redisRepository)(nil)

func NewRedisRepository(redisClient *redis.Client) *redisRepository {
	return &redisRepository{
		redisClient: redisClient,
	}
}

func (r *redisRepository) CacheProduct(ctx context.Context, product *model.Product) error {
	return r.redisClient.Set(ctx, product.ID.Hex(), product, 0).Err() // Adjust expiration if needed
}

func (r *redisRepository) GetCachedProduct(ctx context.Context, id string) (*model.Product, error) {
	val, err := r.redisClient.Get(ctx, id).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	product, err := unmarshalProduct(val)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (r *redisRepository) DeleteCachedProduct(ctx context.Context, id string) error {
	return r.redisClient.Del(ctx, id).Err()
}

func unmarshalProduct(data string) (*model.Product, error) {
	var product model.Product
	err := json.Unmarshal([]byte(data), &product)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling product: %w", err)
	}
	return &product, nil
}
