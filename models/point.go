package models

import "errors"

type Point struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

var (
	ErrInvalidLat = errors.New("Latitude must be between -90 and 90 degrees")
	ErrInvalidLon = errors.New("Longitude must be between -180 and 180 degrees")
)

func (p *Point) Validate() error {
	if p.Lat < -90 || p.Lat > 90 || p.Lat == 0.0 {
		return ErrInvalidLat
	}

	if p.Lon < -180 || p.Lon > 180 || p.Lat == 0.0 {
		return ErrInvalidLon
	}

	return nil
}

func createPoint(lat, lon float64) (*Point, error) {
	point := &Point{lat, lon}

	if err := point.Validate(); err != nil {
		return nil, err
	}

	return point, nil
}
