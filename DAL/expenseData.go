package DAL

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"

	"../models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	expenseRepo struct {
		collection         *mongo.Collection
		categoryCollection *mongo.Collection
	}
	// ExpenseData main interface
	ExpenseData interface {
		GetLastHistory(userId primitive.ObjectID, count int64) []models.Expense
		GetDataByID(userId primitive.ObjectID, id primitive.ObjectID) models.Expense
		AddExpense(userId primitive.ObjectID, expense *models.Expense) primitive.ObjectID
		GetSummary(userId primitive.ObjectID) []models.CategorySummary
		IsAnyInCategory(userId primitive.ObjectID, categoryId primitive.ObjectID) bool
	}
)

func (repo expenseRepo) GetLastHistory(userId primitive.ObjectID, count int64) []models.Expense {
	ctx, cancFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancFunc()
	var result []models.Expense
	t := time.Now()
	filter := bson.M{
		"addedDate": bson.M{"$gt": time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())},
		"_userId":   userId,
	}
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
		filter = bson.M{"_id": expense.CategoryID}
		err = repo.categoryCollection.FindOne(ctx, filter).Decode(&expense.Category)
		if err != nil {
			fmt.Println(err)
		}
		result = append(result, expense)
	}

	return result
}

func (repo expenseRepo) GetDataByID(userId primitive.ObjectID, id primitive.ObjectID) models.Expense {
	ctx, cancFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancFunc()
	filter := bson.M{"_id": id, "_userId": userId}
	var result models.Expense
	err := repo.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		fmt.Println(err)
	}
	return result
}

func (repo expenseRepo) AddExpense(userId primitive.ObjectID, expense *models.Expense) primitive.ObjectID {
	expense.AddedDate = time.Now()
	expense.UserID = userId
	ctx, cancFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancFunc()
	res, err := repo.collection.InsertOne(ctx, expense)
	if err != nil {
		fmt.Println(err)
	}
	return res.InsertedID.(primitive.ObjectID)
}

func (repo expenseRepo) GetSummary(userId primitive.ObjectID) []models.CategorySummary {
	ctx, cancFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancFunc()
	pipeline := []bson.M{
		{"$match": bson.M{"_userId": userId}},
		{"$group": bson.M{
			"_id": "$_categoryId",
			"sum": bson.M{"$sum": "$amount"},
		}},
	}

	cursor, err := repo.collection.Aggregate(ctx, pipeline)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer cursor.Close(context.Background())
	var results []models.CategorySummary
	for cursor.Next(context.Background()) {
		var doc models.CategorySummary
		err := cursor.Decode(&doc)
		if err != nil {
			log.Fatal(err)
		}
		filter := bson.M{"_id": doc.ID}
		var category models.Category
		err = repo.categoryCollection.FindOne(ctx, filter).Decode(&category)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		doc.CategoryName = category.Name
		results = append(results, doc)
	}
	return results
}

func (repo expenseRepo) IsAnyInCategory(userId primitive.ObjectID, categoryId primitive.ObjectID) bool {
	ctx, cancFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancFunc()
	filter := bson.M{"_categoryId": categoryId, "_userId": userId}
	count, _ := repo.collection.CountDocuments(ctx, filter)
	return count > 0
}
