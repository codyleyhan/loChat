package services

import (
	"github.com/codyleyhan/loChat/stores"
	"github.com/codyleyhan/loChat/utils"
	"github.com/labstack/echo"
)

//RegisterServices adds routes to the application
func RegisterServices(app *echo.Echo, userStore stores.UserStore, jwt *utils.JWTGen) {
	userGroup := app.Group("/api/v1/auth")
	registerUserService(userGroup, userStore, jwt)

	tokenMiddleware := utils.CheckToken(jwt)

	api := app.Group("/api/v1", tokenMiddleware)

	// Example of a secure route
	api.GET("/secure", func(c echo.Context) error {
		user := c.Get("user")

		return c.JSON(200, map[string]interface{}{
			"Hello": "World",
			"user":  user,
		})
	})
}
