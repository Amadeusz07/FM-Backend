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
	connection struct {
		database *mongo.Database
	}
	// ExpenseData main interface
	ExpenseData interface {
		GetDataByID(id primitive.ObjectID) models.Expense
		AddExpense(expense *models.Expense) primitive.ObjectID
	}
)

func (conn connection) GetDataByID(id primitive.ObjectID) models.Expense {
	collection := conn.database.Collection("expenses")
	ctx, cancFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancFunc()
	filter := bson.M{"_id": id}
	var result models.Expense
	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}
	return result
}

func (conn connection) AddExpense(expense *models.Expense) primitive.ObjectID {
	collection := conn.database.Collection("expenses")
	ctx, cancFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancFunc()
	res, err := collection.InsertOne(ctx, expense)
	if err != nil {
		log.Fatal(err)
	}
	return res.InsertedID.(primitive.ObjectID)
}

// NewExpenseData creates new database
func NewExpenseData(env string, client *mongo.Client) ExpenseData {
	switch env {
	case "local":
		return connection{
			database: client.Database("testing"),
		}
	}
	return nil
}
