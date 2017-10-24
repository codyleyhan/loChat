package utils

import (
	"net/http"
	"strings"

	"github.com/fatih/color"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type res map[string]interface{}

func AddMiddleware(app *echo.Echo, logo string) {
	app.HideBanner = true
	app.Debug = true

	color.Green(logo)
	color.Cyan("--------------------------------------------------------------------")

	app.Pre(middleware.RequestID())
	app.Pre(middleware.RemoveTrailingSlash())
	app.Use(middleware.Logger())
	app.Use(middleware.Gzip())
	app.Use(middleware.Secure())
	app.Use(middleware.Recover())
	app.Use(ensureJSONRequest)
}

func CheckToken(jwt *JWTGen) func(echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authorization := c.Request().Header.Get("Authorization")

			if authorization == "" {
				return c.JSON(http.StatusUnauthorized, res{
					"message": "You must be logged in.",
				})
			}

			splitAuth := strings.Split(authorization, " ")

			if len(splitAuth) != 2 || splitAuth[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, res{
					"message": "Header should be Authorization: Bearer {token}",
				})
			}

			claims, err := jwt.DecodeToken(splitAuth[1])
			if err != nil {
				return c.JSON(http.StatusUnauthorized, res{
					"message": "The supplied token is invalid",
					"err":     err.Error(),
				})
			}

			c.Set("user", claims)

			return next(c)
		}
	}
}

func ensureJSONRequest(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		request := c.Request()

		if request.Method != "GET" && request.Method != "DELTE" {
			contentType := c.Request().Header.Get(echo.HeaderContentType)

			if contentType != echo.MIMEApplicationJSON {
				return c.JSON(http.StatusUnsupportedMediaType, res{
					"message": "Content-Type must be application/json",
				})
			}
		}

		return next(c)
	}
}
