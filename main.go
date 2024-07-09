package main

import (
	controller "product-app/controller"
	auth "product-app/jwt"
	"product-app/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Gin router'ını başlat
	router := gin.Default()

	// Food routes tanımla
	routes.FoodRoutes(router)
	routes.ExerciseRoutes(router)
	routes.UserRoutes((router))
	protected := router.Group("/user")

	protected.Use(auth.AuthMiddleware())
	protected.DELETE("/", controller.ProtectedEndpoint, controller.DeleteUser)
	// Sunucuyu başlat"
	router.Run(":8080")
}
