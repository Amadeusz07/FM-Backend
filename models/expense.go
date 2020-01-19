package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Expense model
type Expense struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID     primitive.ObjectID `json:"userId,omitempty" bson:"_userId,omitempty"`
	AccountID  primitive.ObjectID `json:"accountId,omitempty" bson:"_accountId,omitempty"`
	CategoryID primitive.ObjectID `json:"categoryId,omitempty" bson:"_categoryId,omitempty"`
	Amount     float32            `json:"amount,omitempty" bson:"amount,omitempty"`
	AddedDate  time.Time          `json:"addedDate,omitempty" bson:"addedDate,omitempty"`
	Category   Category           `json:"category,omitempty"`
}
