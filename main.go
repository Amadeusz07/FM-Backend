package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const mongoURL = "mongodb://localhost:27017"
const env = "local"

func main() {

	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURL))
	if err != nil {
		panic("Error on creating MongoDB Client")
	}
	ctx, canc := context.WithTimeout(context.Background(), 10*time.Second)
	defer canc()
	err = client.Connect(ctx)
	if err != nil {
		panic("Error on connecting to MongoDB")
	}
	fmt.Printf("Connected to MongoDB server %s \n", mongoURL)

	// database := DAL.NewDatabase(env, client)

	// expenseData := DAL.NewExpenseData(database.GetCollection("expenses"))

	// databaseExpense := DAL.NewExpenseData("local", client)
	// id, _ := primitive.ObjectIDFromHex("5e24ab04c9d19abe12b04c0a")
	// result := databaseExpense.GetDataByID(id)

	// // for i := 0; i < 10; i++ {
	// // 	UserID, _ := primitive.ObjectIDFromHex("1")
	// // 	AccountID, _ := primitive.ObjectIDFromHex("2")
	// // 	CategoryID, _ := primitive.ObjectIDFromHex("3")
	// // 	expense := &models.Expense{
	// // 		Amount:     float32(i + 2),
	// // 		AddedDate:  time.Now(),
	// // 		UserID:     UserID,
	// // 		AccountID:  AccountID,
	// // 		CategoryID: CategoryID,
	// // 	}
	// // 	databaseExpense.AddExpense(expense)
	// // }

	// // UserID, _ := primitive.ObjectIDFromHex("1")
	// // AccountID, _ := primitive.ObjectIDFromHex("2")
	// // CategoryID, _ := primitive.ObjectIDFromHex("3")
	// expense := &models.Expense{
	// 	Amount:     float32(66),
	// 	AddedDate:  time.Now(),
	// 	UserID:     primitive.NewObjectID(),
	// 	AccountID:  primitive.NewObjectID(),
	// 	CategoryID: primitive.NewObjectID(),
	// }
	// databaseExpense.AddExpense(expense)

	// fmt.Printf("Found result %f %d", result.Amount, result.AddedDate.Year())
}
