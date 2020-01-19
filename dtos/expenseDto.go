package dtos

import (
	"time"

	"../models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ExpenseDto struct {
	ID         primitive.ObjectID `json:"_id,omitempty"`
	UserID     primitive.ObjectID `json:"userId,omitempty"`
	AccountID  primitive.ObjectID `json:"accountId,omitempty"`
	CategoryID primitive.ObjectID `json:"categoryId,omitempty"`
	Amount     float32            `json:"amount,omitempty"`
	AddedDate  time.Time          `json:"addedDate,omitempty"`
	Category   models.Category    `json:"category,omitempty"`
}
