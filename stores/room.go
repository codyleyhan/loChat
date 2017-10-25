package stores

import (
	"errors"

	"github.com/codyleyhan/loChat/models"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

type (
	RoomStore interface {
		GetAll() (*[]*models.Room, error)
		Get(id int64) (*models.Room, error)
		GetWithinRadius(email string) (*[]*models.Room, error)
		Create(*models.Room) error
	}

	roomStore struct {
		DB *gorm.DB
	}
)

var (
	ErrNoRoom        = errors.New("Room: no room provided")
	ErrDuplicateName = errors.New("Room: that room name is already taken")
)

func (store *roomStore) GetAll() (*[]*models.Room, error) {
	var rooms []*models.Room

	if err := store.DB.Find(&rooms).Error; err != nil {
		return nil, err
	}

	return &rooms, nil
}

func (store *roomStore) GetWithinRadius(email string) (*[]*models.Room, error) {
	return nil, nil
}

func (store *roomStore) Get(id int64) (*models.Room, error) {
	var room models.Room

	if err := store.DB.First(&room, id).Error; err != nil {
		return nil, err
	}

	return &room, nil
}

func (store *roomStore) Create(r *models.Room) error {
	if r == nil {
		return ErrNoRoom
	}

	if err := store.DB.Create(r).Error; err != nil {
		pqErr := err.(*pq.Error)

		if pqErr.Code.Name() == "unique_violation" {
			return ErrDuplicateName
		}

		return err
	}

	return nil
}

//CreateRoomStore creates the db table and holds a reference to the db connection
func CreateRoomStore(db *gorm.DB) RoomStore {
	db.AutoMigrate(&models.Room{})

	return &roomStore{
		DB: db,
	}
}
