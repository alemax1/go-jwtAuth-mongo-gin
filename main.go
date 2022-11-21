package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
	"test-go-project/routes"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	router := gin.New()
	router.Use(gin.Logger())

	routes.AuthRoutes(router)
	routes.UserRoutes(router)

	router.GET("api/v1", func(context *gin.Context) {
		context.JSON(200, gin.H{"success": "Access granted for api/v1!"})
	})

	router.GET("api/v2", func(context *gin.Context) {
		context.JSON(200, gin.H{"success": "Access granted for api/v2!"})
	})

	err = router.Run("127.0.0.1:" + "8000")
	if err != nil {
		log.Panic("error")
	}
}
