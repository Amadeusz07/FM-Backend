package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Category struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID    primitive.ObjectID `json:"userId,omitempty" bson:"userId,omitempty"`
	AccountID primitive.ObjectID `json:"accountId,omitempty" bson:"accountId,omitempty"`
	Name      string             `json:"name" bson:"name,omitempty"`
}
