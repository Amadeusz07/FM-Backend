package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type CategorySummary struct {
	ID           primitive.ObjectID `json:"categoryId,omitempty" bson:"_id,omitempty"`
	CategoryName string             `json:"categoryName,omitempty"`
	Sum          float64            `json:"sum,omitempty"`
}
