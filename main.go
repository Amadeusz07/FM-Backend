package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"./DAL"
	"./controllers"
	authService "./services/auth"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const mongoURL = "mongodb://localhost:27017"
const env = "local"
const port = ":8080"
const frontendAllowedCORS = "http://localhost:4200"

func main() {
	db := initDbConnection()
	controllers.NewExpensesController(db.NewExpenseData())
	controllers.NewCategoriesController(db.NewCategoryData())
	controllers.NewAuthController(db.NewUserData())
	initHTTPServer()
}

func initHTTPServer() {
	fmt.Println("Starting HTTP Server")
	r := mux.NewRouter()
	r.HandleFunc("/register", controllers.Register).Methods(http.MethodPost)
	r.HandleFunc("/login", controllers.Login).Methods(http.MethodPost)
	r.HandleFunc("/logout", controllers.Logout).Methods(http.MethodPost)
	s := r.PathPrefix("").Subrouter()
	s.Use(authService.IsAuthorized)
	s.HandleFunc("/expenses", controllers.GetExpenses).Queries("count", "{count}").Methods(http.MethodGet)
	s.HandleFunc("/expenses", controllers.AddExpense).Methods(http.MethodPost)
	s.HandleFunc("/expenses/{id}", controllers.GetExpense).Methods(http.MethodGet)

	s.HandleFunc("/categories/category-summary", controllers.GetCategoriesSummary).Methods(http.MethodGet)
	s.HandleFunc("/categories", controllers.GetCategories).Methods(http.MethodGet)
	s.HandleFunc("/categories", controllers.AddCategory).Methods(http.MethodPost)
	s.HandleFunc("/categories/{id}", controllers.GetCategory).Methods(http.MethodGet)
	s.HandleFunc("/categories/{id}", controllers.UpdateCategory).Methods(http.MethodPut)
	s.HandleFunc("/categories/{id}", controllers.DeleteCategory).Methods(http.MethodDelete)

	headersOk := handlers.AllowedHeaders([]string{"Origin", "Content-Type", "X-Auth-Token", "Token", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{frontendAllowedCORS})
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"})

	fmt.Printf("Listen and Serve HTTP server on port %s\n", port)
	http.ListenAndServe(port, handlers.CORS(originsOk, headersOk, methodsOk)(r))
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
