package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type ProjectUser struct {
	UserId    primitive.ObjectID `json:"userId" bson:"userId"`
	ProjectId primitive.ObjectID `json:"projectId" bson:"projectId"`
}
