package DAL

import (
	"context"
	"log"
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
		GetDataByID(id primitive.ObjectID) models.Category
		AddCategory(expense *models.Category) primitive.ObjectID
	}
)

func (repo categoryRepo) GetDataByID(id primitive.ObjectID) models.Category {
	ctx, cancFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancFunc()
	filter := bson.M{"_id": id}
	var result models.Category
	err := repo.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}
	return result
}

func (repo categoryRepo) AddCategory(expense *models.Category) primitive.ObjectID {
	ctx, cancFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancFunc()
	res, err := repo.collection.InsertOne(ctx, expense)
	if err != nil {
		log.Fatal(err)
	}
	return res.InsertedID.(primitive.ObjectID)
}
