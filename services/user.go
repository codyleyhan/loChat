package services

import (
	"net/http"

	"github.com/codyleyhan/loChat/models"
	"github.com/codyleyhan/loChat/stores"
	"github.com/codyleyhan/loChat/utils"
	"github.com/labstack/echo"
)

type userService struct {
	store stores.UserStore
	jwt   *utils.JWTGen
}

func (u *userService) register(c context) error {
	var postedUser models.PostedUser

	if err := c.Bind(&postedUser); err != nil {

		return c.JSON(http.StatusBadRequest, res{
			"message": err.Error(),
		})
	}

	if err := postedUser.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, res{
			"message": err.Error(),
		})
	}

	if err := postedUser.HashPassword(); err != nil {
		return c.JSON(http.StatusInternalServerError, res{
			"message": err.Error(),
		})
	}

	user := models.User{
		Email:    postedUser.Email,
		Password: postedUser.Password,
		Admin:    false,
	}

	if err := u.store.Create(&user); err != nil {
		if err == stores.ErrDuplicateEmail {
			return c.JSON(http.StatusBadRequest, res{
				"message": err.Error(),
			})
		}

		return c.JSON(http.StatusInternalServerError, res{
			"message": err.Error(),
		})
	}

	claims := map[string]interface{}{
		"admin": user.Admin,
		"email": user.Email,
		"id":    user.ID,
	}

	token, err := u.jwt.CreateToken(claims)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, res{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, res{
		"token": token,
	})
}

func (u *userService) login(c context) error {
	var postedUser models.PostedUser

	if err := c.Bind(&postedUser); err != nil {
		return c.JSON(http.StatusBadRequest, res{
			"message": err.Error(),
		})
	}

	user, err := u.store.GetByEmail(postedUser.Email)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, res{
			"message": err.Error(),
		})
	}

	if correct := user.CheckPassword(postedUser.Password); !correct {
		return c.JSON(http.StatusUnauthorized, res{
			"message": stores.ErrLoginProblem.Error(),
		})
	}

	claims := map[string]interface{}{
		"admin": user.Admin,
		"email": user.Email,
		"id":    user.ID,
	}

	token, err := u.jwt.CreateToken(claims)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, res{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, res{
		"token": token,
	})
}

func registerUserService(router *echo.Group, store stores.UserStore, jwt *utils.JWTGen) {
	service := userService{store, jwt}

	router.POST("/register", service.register)
	router.POST("/login", service.login)
}
