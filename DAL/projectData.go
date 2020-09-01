package DAL

import (
	"../models"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	projectRepo struct {
		collection            *mongo.Collection
		userProjectCollection *mongo.Collection
		userCollection        *mongo.Collection
	}

	ProjectData interface {
		AddProject(ownerId primitive.ObjectID, project *models.Project) primitive.ObjectID
		GetProject(id primitive.ObjectID) models.Project
		GetAssignedProjects(userId primitive.ObjectID) []models.Project
		GetProjectsForOwner(ownerId primitive.ObjectID) []models.Project
		UpdateProject(ownerId primitive.ObjectID, id primitive.ObjectID, model *models.Project) models.Project
		DisableProject(projectId primitive.ObjectID)
		AssignUser(projectId primitive.ObjectID, userId primitive.ObjectID)
		UnAssignUser(projectId primitive.ObjectID, userId primitive.ObjectID)
		IsUserAssignedToProject(projectId primitive.ObjectID, userId primitive.ObjectID) bool
		IsOwnerOfProject(projectId primitive.ObjectID, userId primitive.ObjectID) bool
	}
)

func (repo projectRepo) AddProject(ownerId primitive.ObjectID, project *models.Project) primitive.ObjectID {
	ctx, cancFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancFunc()
	project.AssignedUsers = []models.UserDto{}
	project.OwnerId = ownerId
	project.AddedDate = time.Now()
	project.Disabled = false
	res, err := repo.collection.InsertOne(ctx, project)
	if err != nil {
		fmt.Println(err)
	}

	return res.InsertedID.(primitive.ObjectID)
}

func (repo projectRepo) GetProject(id primitive.ObjectID) models.Project {
	ctx, cancFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancFunc()

	var result models.Project
	filter := bson.M{"_id": id}
	repo.collection.FindOne(ctx, filter).Decode(&result)

	return result
}

func (repo projectRepo) GetAssignedProjects(userId primitive.ObjectID) []models.Project {
	ctx, cancFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancFunc()

	var projectUsers []models.ProjectUser
	filter := bson.M{"userId": userId}
	cursor, err := repo.userProjectCollection.Find(ctx, filter)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	for cursor.Next(ctx) {
		var projectBson bson.M
		var project models.ProjectUser
		if err = cursor.Decode(&projectBson); err != nil {
			fmt.Println(err)
		}
		bsonBytes, _ := bson.Marshal(projectBson)
		bson.Unmarshal(bsonBytes, &project)
		projectUsers = append(projectUsers, project)
	}
	var projectsIds []primitive.ObjectID
	for _, project := range projectUsers {
		projectsIds = append(projectsIds, project.ProjectId)
	}

	ctx, cancFunc = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancFunc()
	var result []models.Project
	var resultWithUsers []models.Project

	if projectsIds != nil {
		filter = bson.M{"_id": bson.M{"$in": projectsIds}, "disabled": false}
		cursor, err = repo.collection.Find(ctx, filter)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		for cursor.Next(ctx) {
			var projectBson bson.M
			var project models.Project
			if err = cursor.Decode(&projectBson); err != nil {
				fmt.Println(err)
			}
			bsonBytes, _ := bson.Marshal(projectBson)
			bson.Unmarshal(bsonBytes, &project)
			result = append(result, project)
		}
		for _, r := range result {
			r.AssignedUsers = repo.getAssignedUsers(r.ID)
			resultWithUsers = append(resultWithUsers, r)
		}
	}

	return resultWithUsers
}

func (repo projectRepo) GetProjectsForOwner(ownerId primitive.ObjectID) []models.Project {
	ctx, cancFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancFunc()

	var result []models.Project
	filter := bson.M{"_ownerId": ownerId}
	cursor, err := repo.collection.Find(ctx, filter)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var projectBson bson.M
		var project models.Project
		if err = cursor.Decode(&projectBson); err != nil {
			fmt.Println(err)
		}
		bsonBytes, _ := bson.Marshal(projectBson)
		bson.Unmarshal(bsonBytes, &project)
		result = append(result, project)
	}
	var resultWithUsers []models.Project
	for _, r := range result {
		r.AssignedUsers = repo.getAssignedUsers(r.ID)
		resultWithUsers = append(resultWithUsers, r)
	}

	return resultWithUsers
}

func (repo projectRepo) UpdateProject(ownerId primitive.ObjectID, id primitive.ObjectID, model *models.Project) models.Project {
	ctx, cancFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancFunc()
	filter := bson.M{"_id": id, "_ownerId": ownerId}
	update := bson.M{"$set": bson.M{"name": model.Name}}
	_, err := repo.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Println(err)
	}
	result, err := repo.collection.Find(ctx, filter)
	var project models.Project
	if err = result.Decode(&project); err != nil {
		fmt.Println(err)
	}

	return project
}

func (repo projectRepo) DisableProject(projectId primitive.ObjectID) {
	ctx, cancFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancFunc()
	filter := bson.M{"_id": projectId}
	update := bson.M{"$set": bson.M{"disabled": true, "disabledDate": time.Now()}}
	_, err := repo.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Println(err)
	}
}

func (repo projectRepo) AssignUser(projectId primitive.ObjectID, userId primitive.ObjectID) {
	ctx, cancFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancFunc()

	_, err := repo.userProjectCollection.InsertOne(ctx, bson.M{"projectId": projectId, "userId": userId})
	if err != nil {
		fmt.Println(err)
	}
}

func (repo projectRepo) UnAssignUser(projectId primitive.ObjectID, userId primitive.ObjectID) {
	ctx, cancFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancFunc()

	filter := bson.M{"projectId": projectId, "userId": userId}
	_, err := repo.userProjectCollection.DeleteOne(ctx, filter)
	if err != nil {
		fmt.Println(err)
	}
}

func (repo projectRepo) IsUserAssignedToProject(projectId primitive.ObjectID, userId primitive.ObjectID) bool {
	ctx, cancFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancFunc()

	var result models.ProjectUser
	filter := bson.M{"projectId": projectId, "userId": userId}

	err := repo.userProjectCollection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func (repo projectRepo) IsOwnerOfProject(projectId primitive.ObjectID, userId primitive.ObjectID) bool {
	ctx, cancFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancFunc()

	var result models.ProjectUser
	filter := bson.M{"_id": projectId, "_ownerId": userId}
	err := repo.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func (repo projectRepo) getAssignedUsers(projectId primitive.ObjectID) []models.UserDto {
	ctx, cancFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancFunc()

	filter := bson.M{"projectId": projectId}
	cursor, err := repo.userProjectCollection.Find(ctx, filter)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	var userIds []primitive.ObjectID
	for cursor.Next(ctx) {
		var projectBson bson.M
		var projectUser models.ProjectUser
		if err = cursor.Decode(&projectBson); err != nil {
			fmt.Println(err)
		}
		bsonBytes, _ := bson.Marshal(projectBson)
		bson.Unmarshal(bsonBytes, &projectUser)
		userIds = append(userIds, projectUser.UserId)
	}
	var result []models.UserDto
	if userIds != nil {
		filter = bson.M{"_id": bson.M{"$in": userIds}}
		cursor, err = repo.userCollection.Find(ctx, filter)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		for cursor.Next(ctx) {
			var userBson bson.M
			var user models.UserDto
			if err = cursor.Decode(&userBson); err != nil {
				fmt.Println(err)
			}
			bsonBytes, _ := bson.Marshal(userBson)
			bson.Unmarshal(bsonBytes, &user)
			result = append(result, user)
		}
	}

	return result
}
