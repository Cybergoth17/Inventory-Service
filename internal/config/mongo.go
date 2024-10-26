package config

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type MongoConfig struct {
	Host           string `env:"HOST" envDefault:"localhost"`
	Port           string `env:"PORT" envDefault:"27017"`
	User           string `env:"USER" envDefault:"root"`
	Pass           string `env:"PASS" envDefault:"root"`
	DbName         string `env:"DB_NAME" envDefault:"test"`
	CollectionName string `env:"COLLECTION_NAME" envDefault:"test"`
}

var ProductCollection *mongo.Collection

func InitMongoDB(uri, dbName, collectionName string) *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Failed to create Mongo client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	ProductCollection = client.Database(dbName).Collection(collectionName)
	return client
}
