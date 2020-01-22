package DAL

import (
	"context"
	"fmt"
	"time"

	"../models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	categoryRepo struct {
		collection *mongo.Collection
	}
	// CategoryData main interface
	CategoryData interface {
		GetDataByID(id primitive.ObjectID) models.Category
		GetCategories() []models.Category
		AddCategory(expense *models.Category) primitive.ObjectID
		UpdateCategory(id primitive.ObjectID, category *models.Category)
		DeleteCategory(id primitive.ObjectID)
	}
)

func (repo categoryRepo) GetCategories() []models.Category {
	ctx, cancFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancFunc()
	var result []models.Category
	opts := options.Find().SetSort(bson.D{{"name", 1}})
	cursor, err := repo.collection.Find(ctx, bson.D{}, opts)
	if err != nil {
		fmt.Println(err)
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

func (repo categoryRepo) GetDataByID(id primitive.ObjectID) models.Category {
	ctx, cancFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancFunc()
	filter := bson.M{"_id": id}
	var result models.Category
	err := repo.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		fmt.Println(err)
	}
	return result
}

func (repo categoryRepo) AddCategory(category *models.Category) primitive.ObjectID {
	ctx, cancFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancFunc()
	res, err := repo.collection.InsertOne(ctx, category)
	if err != nil {
		fmt.Println(err)
	}
	return res.InsertedID.(primitive.ObjectID)
}

func (repo categoryRepo) UpdateCategory(id primitive.ObjectID, category *models.Category) {
	ctx, cancFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancFunc()
	filter := bson.M{"_id": category.ID}
	update := bson.M{"$set": bson.M{"name": category.Name}}
	_, err := repo.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Println(err)
	}
}

func (repo categoryRepo) DeleteCategory(id primitive.ObjectID) {
	ctx, cancFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancFunc()
	filter := bson.M{"_id": id}
	_, err := repo.collection.DeleteOne(ctx, filter)
	if err != nil {
		fmt.Println(err)
	}
}
