package stores

import (
	"errors"
	"math"

	"github.com/codyleyhan/loChat/models"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

const earthRadius = 6371

type (
	RoomStore interface {
		GetAll() (*[]*models.Room, error)
		Get(id int64) (*models.Room, error)
		GetWithinRadius(r *models.RoomQuery) (*[]*models.Room, error)
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

func (store *roomStore) GetWithinRadius(r *models.RoomQuery) (*[]*models.Room, error) {
	var rooms []*models.Room

	radius := r.Radius / earthRadius
	maxLat := r.Coordinates.Lat + radToDeg(radius)
	minLat := r.Coordinates.Lat - radToDeg(radius)
	maxLon := r.Coordinates.Lon + radToDeg(math.Asin(radius)/math.Cos(degToRad(r.Coordinates.Lat)))
	minLon := r.Coordinates.Lon - radToDeg(math.Asin(radius)/math.Cos(degToRad(r.Coordinates.Lat)))

	rows, err := store.DB.DB().Query(`
		SELECT *
		FROM rooms 
		WHERE lat > $1 AND lat < $2 
		AND lon > $3 AND lon < $4
	`, minLat, maxLat, minLon, maxLon)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var room models.Room

		if err = rows.Scan(&room.ID, &room.Name, &room.UserID, &room.Lat, &room.Lon, &room.CreatedAt, &room.UpdatedAt); err != nil {
			return nil, err
		}

		rooms = append(rooms, &room)
	}

	return &rooms, nil
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

func radToDeg(num float64) float64 {
	return (num * 180) / math.Pi
}

func degToRad(num float64) float64 {
	return (num * math.Pi) / 180
}
