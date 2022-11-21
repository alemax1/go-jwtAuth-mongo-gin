package routes

import (
	"github.com/gin-gonic/gin"
	controller "test-go-project/controllers"
)

func AuthRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("users/sign-up", controller.SignUp())
	incomingRoutes.POST("users/sign-in", controller.SignIn())
}
