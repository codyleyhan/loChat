package models

import (
	"errors"
	"time"
)

type (
	Room struct {
		ID        int64     `json:"id"`
		Name      string    `json:"name" gorm:"unique_index"`
		UserID    int64     `json:"user_id"`
		Lat       float64   `json:"lat"`
		Lon       float64   `json:"lon"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	PostedRoom struct {
		Name        string `json:"name" gorm:"unique_index"`
		Coordinates Point  `json:"coordinates"`
	}
)

var (
	ErrInvalidName = errors.New("Name must be at least 3 characters")
)

func (r *PostedRoom) Validate() error {
	if len(r.Name) < 3 {
		return ErrInvalidName
	}

	if err := r.Coordinates.Validate(); err != nil {
		return err
	}

	return nil
}
