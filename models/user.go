package models

import (
	"errors"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type (
	User struct {
		ID        int64     `json:"id"`
		Email     string    `json:"email" gorm:"unique_index"`
		Password  string    `json:"-"`
		Admin     bool      `json:"-"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	PostedUser struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	UserClaims struct {
		ID    int64
		Email string
		Admin bool
	}
)

var (
	errPasswordLength = errors.New("Password must be at least 6 characters")
	errEmail          = errors.New("That is not a valid email")

	validEmail = regexp.MustCompile(`^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$`)
)

func (u *PostedUser) HashPassword() error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(bytes)

	return nil
}

func (u *PostedUser) Validate() error {
	if u.Email == "" || !validEmail.MatchString(u.Email) {
		return errEmail
	}

	if len(u.Password) < 6 {
		return errPasswordLength
	}

	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
