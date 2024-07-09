package routes

import (
	controllers "product-app/controller"

	"github.com/gin-gonic/gin"
)

func FoodRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/foods", controllers.CreateFood)
	incomingRoutes.GET("/foods", controllers.GetFoods)
	incomingRoutes.GET("/foods/:id", controllers.GetFoodByID)
	incomingRoutes.PUT("/foods/:id", controllers.UpdateFood)
	incomingRoutes.DELETE("/foods/:id", controllers.DeleteFood)
}
