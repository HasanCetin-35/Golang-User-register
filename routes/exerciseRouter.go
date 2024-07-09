package routes

import (
	controllers "product-app/controller"

	"github.com/gin-gonic/gin"
)

func ExerciseRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/exercise", controllers.CreateExercise)
	incomingRoutes.GET("/exercise", controllers.GetExercises)
	incomingRoutes.GET("/exercise/:id", controllers.GetExerciseById)
	incomingRoutes.PUT("/exercise/:id", controllers.UpdateExercise)
	incomingRoutes.DELETE("/exercise/:id", controllers.DeleteExercise)
}