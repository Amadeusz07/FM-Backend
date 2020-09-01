package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Category struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID    primitive.ObjectID `json:"userId,omitempty" bson:"_userId,omitempty"`
	ProjectID primitive.ObjectID `json:"projectId,omitempty" bson:"projectId,omitempty"`
	Name      string             `json:"name" bson:"name,omitempty"`
}
