package inventory

import (
	"context"
	model "inventory-service/internal/model/inventory"
	"inventory-service/internal/repository/inventory"
)

type Client struct {
	MongoRepo inventory.MongoRepository
	RedisRepo inventory.RedisRepository
}

func NewClient(mongoRepo inventory.MongoRepository, redisRepo inventory.RedisRepository) (Client, error) {
	return Client{
		MongoRepo: mongoRepo,
		RedisRepo: redisRepo,
	}, nil
}

func (c *Client) AddProduct(ctx context.Context, product *model.Product) error {
	err := c.MongoRepo.AddProduct(ctx, product)
	if err != nil {
		return err
	}

	return c.RedisRepo.CacheProduct(ctx, product)
}

func (c *Client) GetProduct(ctx context.Context, id string) (*model.Product, error) {
	product, err := c.RedisRepo.GetCachedProduct(ctx, id)
	if err != nil {
		return nil, err
	}

	if product == nil {
		product, err = c.MongoRepo.GetProduct(ctx, id)
		if err != nil {
			return nil, err
		}

		if product != nil {
			err = c.RedisRepo.CacheProduct(ctx, product)
		}
	}

	return product, nil
}

func (c *Client) UpdateProduct(ctx context.Context, id string, updatedProduct *model.Product) error {
	err := c.MongoRepo.UpdateProduct(ctx, id, updatedProduct)
	if err != nil {
		return err
	}

	return c.RedisRepo.CacheProduct(ctx, updatedProduct)
}

func (c *Client) DeleteProduct(ctx context.Context, id string) error {
	err := c.MongoRepo.DeleteProduct(ctx, id)
	if err != nil {
		return err
	}

	return c.RedisRepo.DeleteCachedProduct(ctx, id)
}
