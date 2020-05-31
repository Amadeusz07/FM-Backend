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
		GetDataByID(userId primitive.ObjectID, id primitive.ObjectID) models.Category
		GetCategories(userId primitive.ObjectID) []models.Category
		AddCategory(userId primitive.ObjectID, expense *models.Category) primitive.ObjectID
		UpdateCategory(userId primitive.ObjectID, id primitive.ObjectID, category *models.Category)
		DeleteCategory(userId primitive.ObjectID, id primitive.ObjectID)
	}
)

func (repo categoryRepo) GetCategories(userId primitive.ObjectID) []models.Category {
	ctx, cancFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancFunc()
	var result []models.Category
	filter := bson.M{"_userId": userId}
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

func (repo categoryRepo) GetDataByID(userId primitive.ObjectID, id primitive.ObjectID) models.Category {
	ctx, cancFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancFunc()
	filter := bson.M{"_id": id, "_userId": userId}
	var result models.Category
	err := repo.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		fmt.Println(err)
	}
	return result
}

func (repo categoryRepo) AddCategory(userId primitive.ObjectID, category *models.Category) primitive.ObjectID {
	ctx, cancFunc := context.WithTimeout(context.Background(), 5*time.Second)
	category.UserID = userId
	defer cancFunc()
	res, err := repo.collection.InsertOne(ctx, category)
	if err != nil {
		fmt.Println(err)
	}
	return res.InsertedID.(primitive.ObjectID)
}

func (repo categoryRepo) UpdateCategory(userId primitive.ObjectID, id primitive.ObjectID, category *models.Category) {
	ctx, cancFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancFunc()
	filter := bson.M{"_id": id, "_userId": userId}
	update := bson.M{"$set": bson.M{"name": category.Name}}
	_, err := repo.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Println(err)
	}
}

func (repo categoryRepo) DeleteCategory(userId primitive.ObjectID, id primitive.ObjectID) {
	ctx, cancFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancFunc()
	filter := bson.M{"_id": id, "_userId": userId}
	_, err := repo.collection.DeleteOne(ctx, filter)
	if err != nil {
		fmt.Println(err)
	}
}
