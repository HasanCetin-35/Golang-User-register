package routes

import (
	controllers "product-app/controller"

	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/user", controllers.SignUp)
	incomingRoutes.GET("/user", controllers.Login)
	incomingRoutes.DELETE("/user/:id", controllers.DeleteUser)
	incomingRoutes.GET("/user/deneme", controllers.Deneme)
}
