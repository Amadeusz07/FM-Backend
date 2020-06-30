package DAL

import (
	"../models"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type (
	userRepo struct {
		collection *mongo.Collection
	}

	UserData interface {
		CreateUser(user *models.User) error
		GetUserByEmail(email string) (models.User, error)
	}
)

func (repo userRepo) GetUserByEmail(email string) (models.User, error) {
	ctx, cancFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancFunc()
	var result models.User
	filter := bson.M{"email": email}
	err := repo.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		fmt.Println(err)
		return models.User{}, err
	}
	return result, nil
}

func (repo userRepo) CreateUser(user *models.User) error {
	ctx, cancFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancFunc()
	var result models.User
	filter := bson.M{"email": user.Email}
	err := repo.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		user.Username = user.Email
		user.CreatedDate = time.Now()
		_, err = repo.collection.InsertOne(ctx, user)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		return errors.New("User with provided email exists")
	}
	return nil
}
