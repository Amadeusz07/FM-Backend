package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"./DAL"
	"./config"
	"./controllers"
	authService "./services/auth"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const frontendAllowedCORS = "http://localhost:4200"

var cfg *config.Configuration

func main() {
	cfg = config.GetConfig()
	db := initDbConnection()
	controllers.NewProjectsController(db.NewProjectData())
	controllers.NewExpensesController(db.NewExpenseData())
	controllers.NewCategoriesController(db.NewCategoryData())
	controllers.NewAuthController(db.NewUserData(), db.NewProjectData())
	initHTTPServer()
}

func initHTTPServer() {
	fmt.Println("Starting HTTP Server")
	r := mux.NewRouter()
	r.HandleFunc("/", controllers.Healthcheck).Methods(http.MethodGet)
	r.HandleFunc("/ping", controllers.Ping).Methods(http.MethodGet)
	r.HandleFunc("/register", controllers.Register).Methods(http.MethodPost)
	r.HandleFunc("/login", controllers.Login).Methods(http.MethodPost)
	r.HandleFunc("/logout", controllers.Logout).Methods(http.MethodPost)

	s := r.PathPrefix("").Subrouter()
	s.Use(authService.IsAuthorized)

	s.HandleFunc("/projects", controllers.GetProjectByOwnerId).Methods(http.MethodGet)
	s.HandleFunc("/projects/assigned", controllers.GetAssignedProjects).Methods(http.MethodGet)
	s.HandleFunc("/projects", controllers.CreateProject).Methods(http.MethodPost)

	s.HandleFunc("/selectProject", controllers.SelectProject).Methods(http.MethodPost)
	//if is owner
	s.HandleFunc("/projects/{id}", controllers.UpdateProject).Methods(http.MethodPut)
	s.HandleFunc("/projects/{id}/assignUser", controllers.AssignUser).Methods(http.MethodPost)
	s.HandleFunc("/projects/{id}/unAssignUser", controllers.UnAssignUser).Methods(http.MethodPost)
	s.HandleFunc("/projects/{id}", controllers.DisableProject).Methods(http.MethodDelete)

	s.HandleFunc("/expenses", controllers.GetExpenses).Queries("count", "{count}").Queries("date", "{date}").Methods(http.MethodGet)
	s.HandleFunc("/expenses", controllers.AddExpense).Methods(http.MethodPost)
	s.HandleFunc("/expenses/{id}", controllers.GetExpense).Methods(http.MethodGet)
	s.HandleFunc("/expenses/{id}", controllers.DeleteExpense).Methods(http.MethodDelete)

	s.HandleFunc("/categories/category-summary", controllers.GetCategoriesSummary).Methods(http.MethodGet)
	s.HandleFunc("/categories", controllers.GetCategories).Methods(http.MethodGet)
	s.HandleFunc("/categories", controllers.AddCategory).Methods(http.MethodPost)
	s.HandleFunc("/categories/{id}", controllers.GetCategory).Methods(http.MethodGet)
	s.HandleFunc("/categories/{id}", controllers.UpdateCategory).Methods(http.MethodPut)
	s.HandleFunc("/categories/{id}", controllers.DeleteCategory).Methods(http.MethodDelete)

	headersOk := handlers.AllowedHeaders([]string{"Origin", "Content-Type", "X-Auth-Token", "Token", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{frontendAllowedCORS})
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"})
	fmt.Printf("Listen and Serve HTTP server on port %s\n", cfg.Port)
	http.ListenAndServe(cfg.Port, handlers.CORS(originsOk, headersOk, methodsOk)(r))
}

func initDbConnection() DAL.Database {
	fmt.Println("Starting connection to MongoDB")
	client, err := mongo.NewClient(options.Client().ApplyURI(cfg.ConnectionString).SetDirect(true))
	if err != nil {
		panic("Error on creating MongoDB Client")
	}
	ctx, canc := context.WithTimeout(context.Background(), 10*time.Second)
	defer canc()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal("Error on connecting to MongoDB")
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Connected to MongoDB server %s \n", cfg.ConnectionString)

	return DAL.NewDatabase(cfg, client)
}
