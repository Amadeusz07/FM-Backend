package DAL

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"

	"../models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	expenseRepo struct {
		collection *mongo.Collection
	}
	// ExpenseData main interface
	ExpenseData interface {
		GetLastHistory(count int64) []models.Expense
		GetDataByID(id primitive.ObjectID) models.Expense
		AddExpense(expense *models.Expense) primitive.ObjectID
	}
)

func (repo expenseRepo) GetLastHistory(count int64) []models.Expense {
	ctx, cancFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancFunc()
	var result []models.Expense
	t := time.Now()
	filter := bson.M{"addedDate": bson.M{"$gt": time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())}}
	opts := options.Find().SetSort(bson.D{{"addedDate", -1}}).SetLimit(count)
	cursor, err := repo.collection.Find(ctx, filter, opts)
	if err != nil {
		fmt.Println(err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var expenseBson bson.M
		var expense models.Expense
		if err = cursor.Decode(&expenseBson); err != nil {
			fmt.Println(err)
		}
		bsonBytes, _ := bson.Marshal(expenseBson)
		bson.Unmarshal(bsonBytes, &expense)
		result = append(result, expense)
	}

	return result
}

func (repo expenseRepo) GetDataByID(id primitive.ObjectID) models.Expense {
	ctx, cancFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancFunc()
	filter := bson.M{"_id": id}
	var result models.Expense
	err := repo.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		fmt.Println(err)
	}
	return result
}

func (repo expenseRepo) AddExpense(expense *models.Expense) primitive.ObjectID {
	expense.AddedDate = time.Now()
	ctx, cancFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancFunc()
	res, err := repo.collection.InsertOne(ctx, expense)
	if err != nil {
		fmt.Println(err)
	}
	return res.InsertedID.(primitive.ObjectID)
}
