package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Project struct {
	ID            primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	OwnerId       primitive.ObjectID `json:"ownerId,omitempty" bson:"_ownerId,omitempty"`
	Name          string             `json:"name,omitempty" bson:"name,omitempty"`
	AddedDate     time.Time          `json:"addedDate,omitempty" bson:"addedDate,omitempty"`
	DisabledDate  time.Time          `json:"disabledDate,omitempty" bson:"disabledDate"`
	Disabled      bool               `json:"disabled" bson:"disabled"`
	AssignedUsers []User             `json:"assignedUsers" bson:"assignedUsers"`
}
