package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Expense model
type Expense struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID     primitive.ObjectID `json:"userId,omitempty" bson:"_userId,omitempty"`
	ProjectID  primitive.ObjectID `json:"projectId,omitempty" bson:"projectId,omitempty"`
	CategoryID primitive.ObjectID `json:"categoryId,omitempty" bson:"_categoryId,omitempty"`
	Title      string             `json:"title,omitempty" bson:"title,omitempty"`
	Amount     float64            `json:"amount,omitempty" bson:"amount,omitempty"`
	AddedDate  time.Time          `json:"addedDate,omitempty" bson:"addedDate,omitempty"`
	Category   Category           `json:"category,omitempty" bson:"-"`
}
