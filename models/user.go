package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username    string             `json:"username" bson:"username"`
	CreatedDate time.Time          `json:"createdDate" bson:"createdDate"`
}
