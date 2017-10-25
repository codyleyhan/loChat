package services

import (
	"net/http"

	"github.com/codyleyhan/loChat/models"
	"github.com/codyleyhan/loChat/stores"
	"github.com/labstack/echo"
)

type roomService struct {
	store stores.RoomStore
}

func (r *roomService) getAll(c context) error {
	// user := getUser(c)

	// if !user.Admin {
	// 	return c.JSON(http.StatusForbidden, res{
	// 		"message": "You are not allowed.",
	// 	})
	// }

	rooms, err := r.store.GetAll()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, res{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, res{
		"rooms": rooms,
	})
}

func (r *roomService) create(c context) error {
	var postedRoom models.PostedRoom

	if err := c.Bind(&postedRoom); err != nil {
		return c.JSON(http.StatusBadRequest, res{
			"message": err.Error(),
		})
	}

	if err := postedRoom.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, res{
			"message": err.Error(),
		})
	}

	user := getUser(c)

	room := models.Room{
		Name:   postedRoom.Name,
		UserID: user.ID,
		Lat:    postedRoom.Coordinates.Lat,
		Lon:    postedRoom.Coordinates.Lon,
	}

	if err := r.store.Create(&room); err != nil {
		if err == stores.ErrDuplicateName {
			return c.JSON(http.StatusBadRequest, res{
				"message": err.Error(),
			})
		}

		return c.JSON(http.StatusInternalServerError, res{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, res{
		"room": room,
	})
}

func registerRoomService(router *echo.Group, store stores.RoomStore) {
	service := roomService{store}

	router.GET("", service.getAll)
	router.POST("", service.create)
}
