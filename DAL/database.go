package DAL

// import (
// 	"go.mongodb.org/mongo-driver/mongo"
// )

// type (
// 	connection struct {
// 		database *mongo.Database
// 	}

// 	// Database interface
// 	Database interface {
// 		GetCollection(name string) mongo.Collection
// 	}
// )

// func (conn connection) GetCollection(name string) *mongo.Collection {
// 	return conn.database.Collection(name)
// }

// func NewDatabase(env string, client *mongo.Client) Database {
// 	switch env {
// 	case "local":
// 		return connection{
// 			database: client.Database("testing"),
// 		}
// 	}
// 	return nil
// }

// // NewExpenseData creates new database
// func NewExpenseData(env string, client *mongo.Client) ExpenseData {
// 	switch env {
// 	case "local":
// 		return connection{
// 			database: client.Database("testing"),
// 		}
// 	}
// 	return nil
// }
