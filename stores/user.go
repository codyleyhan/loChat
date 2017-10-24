package stores

import (
	"errors"
	"fmt"

	"github.com/codyleyhan/loChat/models"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

type (
	UserStore interface {
		GetAll() (*[]*models.User, error)
		Get(id int64) (*models.User, error)
		GetByEmail(email string) (*models.User, error)
		Create(*models.User) error
	}

	userStore struct {
		DB *gorm.DB
	}
)

var (
	ErrNoUser           = errors.New("No user provided")
	ErrDuplicateEmail   = errors.New("That email is already registered with fixMoto")
	ErrGenericDBProblem = errors.New("There was a problem registering user please try again")
	ErrLoginProblem     = errors.New("That email/password is incorrect")
)

func (store *userStore) GetAll() (*[]*models.User, error) {
	var users []*models.User

	if err := store.DB.Find(&users).Error; err != nil {
		return nil, err
	}

	return &users, nil
}

func (store *userStore) Get(id int64) (*models.User, error) {
	var user models.User

	if err := store.DB.First(&user, id).Error; err != nil {
		return nil, ErrLoginProblem
	}

	return &user, nil
}

func (store *userStore) GetByEmail(email string) (*models.User, error) {
	var user models.User

	if err := store.DB.Where("email = ?", email).First(&user).Error; err != nil {
		fmt.Println(err)
		return nil, ErrLoginProblem
	}

	return &user, nil
}

func (store *userStore) Create(u *models.User) error {
	if u == nil {
		return ErrNoUser
	}

	if err := store.DB.Create(u).Error; err != nil {
		pqErr := err.(*pq.Error)

		if pqErr.Code.Name() == "unique_violation" {
			return ErrDuplicateEmail
		}

		return ErrGenericDBProblem
	}

	return nil
}

//CreateUserStore creates the db table and holds a reference to the db connection
func CreateUserStore(db *gorm.DB) UserStore {
	db.AutoMigrate(&models.User{})

	return &userStore{
		DB: db,
	}
}
