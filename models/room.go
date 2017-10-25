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

	RoomQuery struct {
		Coordinates Point   `json:"coordinates"`
		Radius      float64 `json:"radius"`
	}
)

var (
	ErrInvalidName   = errors.New("Name must be at least 3 characters")
	ErrInvalidRadius = errors.New("Radius must be greater than 0 miles")
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

func (r *RoomQuery) Validate() error {
	if err := r.Coordinates.Validate(); err != nil {
		return err
	}

	if r.Radius <= 0 {
		return ErrInvalidRadius
	}

	return nil
}
