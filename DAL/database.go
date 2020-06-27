package DAL

import (
	"go.mongodb.org/mongo-driver/mongo"

	"../config"
)

type (
	connection struct {
		database *mongo.Database
	}

	// Database interface
	Database interface {
		NewExpenseData() ExpenseData
		NewCategoryData() CategoryData
		NewUserData() UserData
		NewProjectData() ProjectData
	}
)

// NewDatabase creates connection to database
func NewDatabase(cfg *config.Configuration, client *mongo.Client) Database {
	return connection{
		database: client.Database(cfg.DatabaseName),
	}
}

func (conn connection) NewExpenseData() ExpenseData {
	return expenseRepo{
		collection:         conn.database.Collection("expenses"),
		categoryCollection: conn.database.Collection("categories"),
	}
}

func (conn connection) NewCategoryData() CategoryData {
	return categoryRepo{
		collection: conn.database.Collection("categories"),
	}
}

func (conn connection) NewUserData() UserData {
	return userRepo{
		collection: conn.database.Collection("users"),
	}
}

func (conn connection) NewProjectData() ProjectData {
	return projectRepo{
		collection:            conn.database.Collection("projects"),
		userProjectCollection: conn.database.Collection("projects_users"),
	}
}
