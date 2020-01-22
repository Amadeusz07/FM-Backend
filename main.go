package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"./DAL"
	"./controllers"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const mongoURL = "mongodb://localhost:27017"
const env = "local"
const port = ":8080"

func main() {
	db := initDbConnection()
	controllers.NewExpensesController(db.NewExpenseData())
	controllers.NewCategoriesController(db.NewCategoryData())
	initHTTPServer()
}

func initHTTPServer() {
	fmt.Println("Starting HTTP Server")
	r := mux.NewRouter()
	r.HandleFunc("/expenses/{id}", controllers.GetExpense).Methods(http.MethodGet)
	r.HandleFunc("/expenses", controllers.AddExpense).Methods(http.MethodPost)

	r.HandleFunc("/categories", controllers.GetCategories).Methods(http.MethodGet)
	r.HandleFunc("/categories", controllers.AddCategory).Methods(http.MethodPost)
	r.HandleFunc("/categories/{id}", controllers.GetCategory).Methods(http.MethodGet)
	r.HandleFunc("/categories/{id}", controllers.UpdateCategory).Methods(http.MethodPut)
	r.HandleFunc("/categories/{id}", controllers.DeleteCategory).Methods(http.MethodDelete)

	r.Use(mux.CORSMethodMiddleware(r))

	fmt.Printf("Listen and Serve HTTP server on port %s\n", port)
	http.ListenAndServe(port, r)

}

func initDbConnection() DAL.Database {
	fmt.Println("Starting connection to MongoDB")
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

	return DAL.NewDatabase("local", client)
}
