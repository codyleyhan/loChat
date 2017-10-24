package utils

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq" // driver for postgres
)

//DbOptions options that will used to connect to db
type DbOptions struct {
	User     string
	Password string
	Name     string
	Port     string
}

//ConnectToDB creates connection to database
func ConnectToDB(opts DbOptions) *gorm.DB {
	dbConnection := fmt.Sprintf("user=%s password=%s dbname=%s port=%s sslmode=disable",
		opts.User, opts.Password, opts.Name, opts.Port)

	var db *gorm.DB
	var err error
	if db, err = gorm.Open("postgres", dbConnection); err != nil {
		log.Fatal("❌  DB CONN: " + err.Error())
	}

	fmt.Println("⚡  DB Connected")

	return db
}
