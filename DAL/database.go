package DAL

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	connection struct {
		database *mongo.Database
	}

	// Database interface
	Database interface {
		NewExpenseData() ExpenseData
		NewCategoryData() CategoryData
	}
)

// NewDatabase creates connection to database
func NewDatabase(env string, client *mongo.Client) Database {
	switch env {
	case "local":
		return connection{
			database: client.Database("testing"),
		}
	}
	return nil
}

func (conn connection) NewExpenseData() ExpenseData {
	return expenseRepo{
		collection: conn.database.Collection("expenses"),
	}
}

func (conn connection) NewCategoryData() CategoryData {
	return categoryRepo{
		collection: conn.database.Collection("categories"),
	}
}
