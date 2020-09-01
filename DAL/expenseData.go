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
		GetLastHistory(projectId primitive.ObjectID, count int64, date time.Time) []models.Expense
		GetDataByID(projectId primitive.ObjectID, id primitive.ObjectID) models.Expense
		AddExpense(userId primitive.ObjectID, projectId primitive.ObjectID, expense *models.Expense) primitive.ObjectID
		DeleteExpense(projectId primitive.ObjectID, expenseId primitive.ObjectID)
		GetSummary(projectId primitive.ObjectID) []models.CategorySummary
		IsAnyInCategory(projectId primitive.ObjectID, categoryId primitive.ObjectID) bool
	}
)

func (repo expenseRepo) GetLastHistory(projectId primitive.ObjectID, count int64, date time.Time) []models.Expense {
	ctx, cancFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancFunc()
	var result []models.Expense
	//t := time.Now()
	fromDate := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.UTC)
	toDate := time.Date(date.Year(), date.Month()+1, 1, 0, 0, 0, 0, time.UTC)
	filter := bson.M{
		"addedDate": bson.M{
			"$gt": fromDate,
			"$lt": toDate,
		},
		"projectId": projectId,
	}
	//opts := options.Find().SetSort(bson.D{{"addedDate", -1}})
	opts := options.Find().SetLimit(count)
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

func (repo expenseRepo) GetDataByID(projectId primitive.ObjectID, id primitive.ObjectID) models.Expense {
	ctx, cancFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancFunc()
	filter := bson.M{"_id": id, "projectId": projectId}
	var result models.Expense
	err := repo.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		fmt.Println(err)
	}
	return result
}

func (repo expenseRepo) AddExpense(userId primitive.ObjectID, projectId primitive.ObjectID, expense *models.Expense) primitive.ObjectID {
	expense.AddedDate = time.Now()
	expense.UserID = userId
	expense.ProjectID = projectId
	ctx, cancFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancFunc()
	res, err := repo.collection.InsertOne(ctx, expense)
	if err != nil {
		fmt.Println(err)
	}
	return res.InsertedID.(primitive.ObjectID)
}

func (repo expenseRepo) DeleteExpense(projectId primitive.ObjectID, id primitive.ObjectID) {
	ctx, cancFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancFunc()
	filter := bson.M{"_id": id, "projectId": projectId}
	_, err := repo.collection.DeleteOne(ctx, filter)
	if err != nil {
		fmt.Println(err)
	}
}

func (repo expenseRepo) GetSummary(projectId primitive.ObjectID) []models.CategorySummary {
	ctx, cancFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancFunc()
	pipeline := []bson.M{
		{"$match": bson.M{"projectId": projectId}},
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

func (repo expenseRepo) IsAnyInCategory(projectId primitive.ObjectID, categoryId primitive.ObjectID) bool {
	ctx, cancFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancFunc()
	filter := bson.M{"_categoryId": categoryId, "projectId": projectId}
	count, _ := repo.collection.CountDocuments(ctx, filter)
	return count > 0
}
