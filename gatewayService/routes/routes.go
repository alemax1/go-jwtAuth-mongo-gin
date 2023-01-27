package routes

import (
	"amaksimov/gatewayService/controllers"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

func GatewayApi() {
	e := echo.New()

	apiGroup := e.Group("/api/v1")

	apiGroup.GET("/user/:id/cars", controllers.GetUserCars, CtxMiddleware)
	apiGroup.GET("/user/:id/engines", controllers.GetUserEngines, CtxMiddleware)
	apiGroup.GET("/car/:id/engine", controllers.GetCarEngine, CtxMiddleware)
	apiGroup.GET("/engines", controllers.GetAllEngines)

	e.Logger.Fatal(e.Start(":" + viper.GetString("routes.port")))
}
