package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/stepup99/go_mongo/controllers"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var userCollection *mongo.Collection

func init() {
	// create clientoption
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	userCollection = client.Database("testdb").Collection("users")

}

func main() {
	// Create a UserCollection instance
	userController := &controllers.UserController{Collection: userCollection}
	// Set up Gin
	r := gin.Default()
	r.GET("/users", userController.GetAllUsers)
	r.POST("/users", userController.CreateUser)
	r.PUT("user/:id", userController.UpdateUser)
	r.DELETE("user/:id", userController.DeleteUser)
	r.Run(":8080")
}
