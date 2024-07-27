package controllers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stepup99/go_mongo/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserController struct {
	Collection *mongo.Collection
}

func (uc *UserController) GetAllUsers(c *gin.Context) {
	var users []models.User
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := uc.Collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching error"})
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding user"})
			return
		}
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cursor error"})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (uc *UserController) CreateUser(c *gin.Context) {
	var user models.User
	// read the request body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed while reading the request body"})
		return
	}

	// unmarshal the request
	if err := json.Unmarshal(body, &user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON format"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// insert the user to the db

	result, err := uc.Collection.InsertOne(ctx, user)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to insert the data"})
		return
	}

	// convert the userID to string
	id, ok := result.InsertedID.(primitive.ObjectID)

	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to convert UserID"})
		return
	}

	user.ID = id.Hex()

	c.JSON(http.StatusCreated, user)

}

func (uc *UserController) DeleteUser(c *gin.Context) {
	// get the ID from param
	id := c.Param("id")

	// Validate the userID
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid ID format"})
	}

	// create context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// delete from database
	result, err := uc.Collection.DeleteOne(ctx, bson.M{"_id": oid})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while deleteing the user"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "User not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User Deleted successfully!!!"})
}

func (uc *UserController) UpdateUser(c *gin.Context) {
	var user models.User

	// get the id
	id := c.Param("id")
	// validate the id
	objid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, "error ID format")
		return
	}

	// get the request body

	body, err := io.ReadAll(c.Request.Body)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while reading the request"})
		return
	}
	// unmarshal the request body
	err = json.Unmarshal(body, &user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body format"})
		return
	}

	// create the context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// update the data
	update := bson.M{
		"$set": user,
	}
	result, err := uc.Collection.UpdateOne(ctx, bson.M{"_id": objid}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating user"})
		return
	}
	// check if data matched or not
	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
		return
	}

	// return the data
	c.JSON(http.StatusOK, gin.H{"message": "updated successfully !!!"})
}
