package services

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/codyleyhan/loChat/models"
	"github.com/codyleyhan/loChat/stores"
	"github.com/codyleyhan/loChat/utils"
	"github.com/labstack/echo"
)

//RegisterServices adds routes to the application
func RegisterServices(app *echo.Echo, userStore stores.UserStore, roomStore stores.RoomStore, jwt *utils.JWTGen) {
	userGroup := app.Group("/api/v1/auth")
	registerUserService(userGroup, userStore, jwt)

	tokenMiddleware := utils.CheckToken(jwt)

	api := app.Group("/api/v1", tokenMiddleware)

	roomsGroup := api.Group("/rooms")
	registerRoomService(roomsGroup, roomStore)

	// Example of a secure route
	api.GET("/secure", func(c echo.Context) error {
		user := c.Get("user")

		return c.JSON(200, map[string]interface{}{
			"Hello": "World",
			"user":  user,
		})
	})
}

func getUser(c context) *models.UserClaims {
	claims := c.Get("user").(map[string]interface{})
	userClaims := models.UserClaims{
		ID:    int64(claims["id"].(float64)),
		Email: claims["email"].(string),
		Admin: claims["admin"].(bool),
	}

	fmt.Println("HERE", userClaims)

	return &userClaims
}

func checkID(next func(int64) echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		potentialID := c.Param("id")
		id, err := strconv.ParseInt(potentialID, 10, 64)

		if err != nil || id <= 0 {
			return echo.NewHTTPError(http.StatusBadRequest, potentialID+" is not a valid id.")
		}

		return next(id)(c)
	}
}
