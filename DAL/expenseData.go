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
	expenseRepo struct {
		collection *mongo.Collection
	}
	// ExpenseData main interface
	ExpenseData interface {
		GetDataByID(id primitive.ObjectID) models.Expense
		AddExpense(expense *models.Expense) primitive.ObjectID
	}
)

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
