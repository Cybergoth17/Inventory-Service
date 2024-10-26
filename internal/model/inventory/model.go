package inventory

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Product struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Category    string             `json:"category" bson:"category"`
	Price       float64            `json:"price" bson:"price"`
	Description string             `json:"description" bson:"description"`
	Quantity    int                `json:"quantity" bson:"quantity"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
}

type MessageType struct {
	Message string `json:"message"`
}
