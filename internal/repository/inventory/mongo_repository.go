package inventory

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	model "inventory-service/internal/model/inventory"
)

type MongoRepository interface {
	AddProduct(ctx context.Context, product *model.Product) error
	GetProduct(ctx context.Context, id string) (*model.Product, error)
	UpdateProduct(ctx context.Context, id string, updatedProduct *model.Product) error
	DeleteProduct(ctx context.Context, id string) error
}

type mongoRepository struct {
	collection *mongo.Collection
}

var _ MongoRepository = (*mongoRepository)(nil)

func NewMongoRepository(mongoCollection *mongo.Collection) *mongoRepository {
	return &mongoRepository{
		collection: mongoCollection,
	}
}

func (r *mongoRepository) AddProduct(ctx context.Context, product *model.Product) error {
	product.ID = primitive.NewObjectID()
	product.CreatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, product)
	return err
}

func (r *mongoRepository) GetProduct(ctx context.Context, id string) (*model.Product, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var product model.Product
	filter := bson.M{"_id": objID}

	err = r.collection.FindOne(ctx, filter).Decode(&product)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *mongoRepository) UpdateProduct(ctx context.Context, id string, updatedProduct *model.Product) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objID}
	update := bson.M{
		"$set": bson.M{
			"name":        updatedProduct.Name,
			"category":    updatedProduct.Category,
			"price":       updatedProduct.Price,
			"description": updatedProduct.Description,
			"quantity":    updatedProduct.Quantity,
			"updated_at":  time.Now(),
		},
	}

	_, err = r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *mongoRepository) DeleteProduct(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objID}
	_, err = r.collection.DeleteOne(ctx, filter)
	return err
}
