package routes

import (
	"github.com/gin-gonic/gin"
	"backend/controllers"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/users/signup", controllers.SignUp())
}
