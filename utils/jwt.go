package utils

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTGen struct {
	NumDaysValid int
	Secret       []byte
}

var ErrInvalidToken = errors.New("That is not a valid token.")

//CreateToken generates a token witht the given claims
func (j *JWTGen) CreateToken(data map[string]interface{}) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().AddDate(0, 0, j.NumDaysValid).Unix() * 1000

	for k, v := range data {
		claims[k] = v
	}

	tokenString, err := token.SignedString(j.Secret)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

//DecodeToken given a token will decode and validate the token
func (j *JWTGen) DecodeToken(tokenString string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return j.Secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}
