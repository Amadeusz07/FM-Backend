package DAL

import (
	"context"
	"fmt"
	"time"

	"../models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	categoryRepo struct {
		collection *mongo.Collection
	}
	// CategoryData main interface
	CategoryData interface {
		GetDataByID(projectId primitive.ObjectID, id primitive.ObjectID) models.Category
		GetCategories(projectId primitive.ObjectID) []models.Category
		AddCategory(userId primitive.ObjectID, projectId primitive.ObjectID, expense *models.Category) primitive.ObjectID
		UpdateCategory(projectId primitive.ObjectID, id primitive.ObjectID, category *models.Category)
		DeleteCategory(projectId primitive.ObjectID, id primitive.ObjectID)
	}
)

func (repo categoryRepo) GetCategories(projectId primitive.ObjectID) []models.Category {
	ctx, cancFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancFunc()
	var result []models.Category
	filter := bson.M{"projectId": projectId}
	//opts := options.Find().SetSort(bson.D{{"name", 1}})
	cursor, err := repo.collection.Find(ctx, filter, nil)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var categoryBson bson.M
		var category models.Category
		if err = cursor.Decode(&categoryBson); err != nil {
			fmt.Println(err)
		}
		bsonBytes, _ := bson.Marshal(categoryBson)
		bson.Unmarshal(bsonBytes, &category)
		result = append(result, category)
	}
	return result
}

func (repo categoryRepo) GetDataByID(projectId primitive.ObjectID, id primitive.ObjectID) models.Category {
	ctx, cancFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancFunc()
	filter := bson.M{"_id": id, "projectId": projectId}
	var result models.Category
	err := repo.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		fmt.Println(err)
	}
	return result
}

func (repo categoryRepo) AddCategory(userId primitive.ObjectID, projectId primitive.ObjectID, category *models.Category) primitive.ObjectID {
	category.ProjectID = projectId
	ctx, cancFunc := context.WithTimeout(context.Background(), 5*time.Second)
	category.UserID = userId
	defer cancFunc()
	res, err := repo.collection.InsertOne(ctx, category)
	if err != nil {
		fmt.Println(err)
	}
	return res.InsertedID.(primitive.ObjectID)
}

func (repo categoryRepo) UpdateCategory(projectId primitive.ObjectID, id primitive.ObjectID, category *models.Category) {
	ctx, cancFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancFunc()
	filter := bson.M{"_id": id, "projectId": projectId}
	update := bson.M{"$set": bson.M{"name": category.Name}}
	_, err := repo.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Println(err)
	}
}

func (repo categoryRepo) DeleteCategory(projectId primitive.ObjectID, id primitive.ObjectID) {
	ctx, cancFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancFunc()
	filter := bson.M{"_id": id, "projectId": projectId}
	_, err := repo.collection.DeleteOne(ctx, filter)
	if err != nil {
		fmt.Println(err)
	}
}
