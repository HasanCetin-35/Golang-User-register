package controllers

import (
	"context"
	"log"
	"net/http"
	"product-app/database"
	"product-app/models"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")

func CreateFood(c *gin.Context) {
	var food models.Food

	if err := c.BindJSON(&food); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uuidWithHyphens := uuid.New().String()
    food.ID = strings.ReplaceAll(uuidWithHyphens, "-", "")
	

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := foodCollection.InsertOne(ctx, food)
	if err != nil {
		log.Fatalf("Error inserting document: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error inserting document"})
		return
	}

	c.JSON(http.StatusCreated, food)
}

func GetFoods(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := foodCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching documents"})
		return
	}
	defer cursor.Close(ctx)

	var foods []models.Food
	for cursor.Next(ctx) {
		var food models.Food
		if err := cursor.Decode(&food); err != nil {
			log.Printf("Error decoding document: %v\n", err)
			continue
		}
		foods = append(foods, food)
	}

	if err := cursor.Err(); err != nil {
		log.Printf("Cursor error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cursor error"})
		return
	}

	c.JSON(http.StatusOK, foods)
}

func GetFoodByID(c *gin.Context) {
	foodID := c.Param("id")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var food models.Food
	err := foodCollection.FindOne(ctx, bson.M{"_id": foodID}).Decode(&food)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		return
	}

	c.JSON(http.StatusOK, food)
}

func UpdateFood(c *gin.Context) {
	foodID := c.Param("id")
	
	var food models.Food
	if err := c.BindJSON(&food); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	update := bson.M{
		"name":  food.Name,
		"price": food.Price,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := foodCollection.UpdateOne(ctx, bson.M{"_id": foodID}, bson.M{"$set": update})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating document"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Document updated successfully"})
}

func DeleteFood(c *gin.Context) {
	foodID := c.Param("id")
	

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := foodCollection.DeleteOne(ctx, bson.M{"_id": foodID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting document"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Document deleted successfully"})
}
