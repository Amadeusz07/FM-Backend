package dtos

import "go.mongodb.org/mongo-driver/bson/primitive"

type ChangeProjectRequest struct {
	ProjectId primitive.ObjectID `json:"projectId"`
}
