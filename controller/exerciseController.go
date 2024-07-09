package controllers

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"product-app/database"
	"product-app/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var exerciseCollection *mongo.Collection = database.OpenCollection(database.Client, "exercise")

func CreateExercise(c *gin.Context) {
	var exercise models.Exercise
	if err := c.BindJSON(&exercise); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uuidWithHyphens := uuid.New().String()
	exercise.ID = strings.ReplaceAll(uuidWithHyphens, "-", "")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := exerciseCollection.InsertOne(ctx, exercise)
	if err != nil {
		log.Fatalf("Error inserting document: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error inserting document"})
		return
	}
	c.JSON(http.StatusCreated, exercise)
}
func GetExercises(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := exerciseCollection.Find(ctx, bson.M{})

	//bson.M{} filtreleme görevi görür bson.Map olarak bakabilirsin.
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching documents"})
		//gin.H json formatında yanıt döndürmek için kullanılır.
	}
	defer cursor.Close(ctx)
	var exercises []models.Exercise
	for cursor.Next(ctx) {
		var exercise models.Exercise
		if err := cursor.Decode(&exercise); err != nil {
			log.Printf("Error decoding document:%v\n", err)
			continue
		}
		exercises = append(exercises, exercise)

	}
	if err := cursor.Err(); err != nil {
		log.Printf("Cursor error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cursor error"})
		return
	}
	c.JSON(http.StatusOK, exercises)
}
func GetExerciseById(c *gin.Context) {
	exerciseId := c.Param("id")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var exercise models.Exercise
	err := exerciseCollection.FindOne(ctx, bson.M{"_id": exerciseId}).Decode(&exercise)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		return
	}
	c.JSON(http.StatusOK, exercise)
}
func UpdateExercise(c *gin.Context) {
	exerciseId := c.Param("id")
	var exercise models.Exercise
	if err := c.BindJSON(&exercise); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	update := bson.M{
		"name":          exercise.Name,
		"exercise_type": exercise.Exercise_Type,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := exerciseCollection.UpdateOne(ctx, bson.M{"_id": exerciseId}, bson.M{"$set": update})
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
func DeleteExercise(c *gin.Context){
	exerciseID := c.Param("id")

	ctx,cancel := context.WithTimeout(context.Background(),10*time.Second)
	defer cancel()

	result,err:= exerciseCollection.DeleteOne(ctx,bson.M{"_id":exerciseID})
	if err !=nil{
		c.JSON(http.StatusInternalServerError,gin.H{"error":"Error deleting document"})
		return
	}
	if result.DeletedCount==0{
		c.JSON(http.StatusFound,gin.H{"error":"Document not found"})
		return
	}
	c.JSON(http.StatusOK,gin.H{"message":"Document deleted successfully"})
}
