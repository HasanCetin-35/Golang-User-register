package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"product-app/database"
	auth "product-app/jwt"
	"product-app/models"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate = validator.New()

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

func SignUp(c *gin.Context) {
	var user models.User
	if err := c.Bind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read user"})
		return
	}
	// E-posta adresinin benzersiz olduğunu kontrol et
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var existingUser models.User
	err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&existingUser)
	if err == nil {
		// Eğer e-posta adresi zaten kullanımda ise hata döndür
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email address already exists"})
		return
	} else if err != mongo.ErrNoDocuments {
		// Veritabanı hatası durumunda
		log.Fatalf("Error checking existing email: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking email availability"})
		return
	}
	// Validasyon kontrolü yapılıyor
	if err := validate.Struct(user); err != nil {
		errorFields := make([]ErrorResponse, 0)
		for _, err := range err.(validator.ValidationErrors) {
			element := ErrorResponse{
				FailedField: err.Field(),
				Tag:         err.Tag(),
				Value:       err.Param(),
			}
			errorFields = append(errorFields, element)
		}
		c.JSON(http.StatusBadRequest, gin.H{"error_fields": errorFields})
		return
	}

	// Şifre hashleme
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = string(hashedPassword)

	uuidWithHyphens := uuid.New().String()
	user.ID = strings.ReplaceAll(uuidWithHyphens, "-", "")

	_, err = userCollection.InsertOne(ctx, user)
	if err != nil {
		log.Fatalf("Error inserting document: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error inserting document"})
		return
	}

	c.JSON(http.StatusCreated, user)
}
func Login(c *gin.Context) {
	var loginUser models.User
	if err := c.BindJSON(&loginUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read user"})
		return
	}
	fmt.Printf("Login User: %+v\n", loginUser)
	var dbUser models.User
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := userCollection.FindOne(ctx, bson.M{"email": loginUser.Email}).Decode(&dbUser)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}
	fmt.Println(bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(loginUser.Password)))
	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(loginUser.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}
	tokenString, err := auth.CreateJWT(dbUser.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
func ProtectedEndpoint(c *gin.Context) {
    email := c.MustGet("email").(string)

    var user models.User
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    err := userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching user data"})
        return
    }

    c.Set("user", &user) 
	
    c.Next() 
}

func Deneme(c *gin.Context) {
    // user'ı models.User tipinde al
    userInterface, exists := c.Get("user")
	fmt.Printf("User ID: %s\n", userInterface)
    if !exists {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found in context"})
        return
    }

    user, ok := userInterface.(*models.User)
    if !ok {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user type in context"})
        return
    }

    // user'ın ID'sini yazdır
    fmt.Printf("User ID: %s\n", user.ID)
	c.JSON(http.StatusOK,gin.H{"user":&user})
    // Diğer user bilgileri ile yapmak istediğiniz işlemleri burada yapabilirsiniz
}

func DeleteUser(c *gin.Context) {
	userInterface, exists := c.Get("user")
	fmt.Printf("User ID: %s\n", userInterface)
    if !exists {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found in context"})
        return
    }
	user, ok := userInterface.(*models.User)
    if !ok {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user type in context"})
        return
    }
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := userCollection.DeleteOne(ctx, bson.M{"_id": user.ID})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting user"})
		return
	}
	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not fount"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

