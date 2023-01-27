package routes

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"
)

func CtxMiddleware(h echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*3)
		defer cancel()

		c.SetRequest(c.Request().WithContext(ctx))

		return h(c)
	}
}
