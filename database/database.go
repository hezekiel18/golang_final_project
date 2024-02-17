package database

import (
	"final_project/models"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	host     = "localhost"
	user     = "farras"
	password = "password123"
	dbPort   = "5433"
	dbName   = "final-project"
	db       *gorm.DB
	err      error
)

func init() {
	host = os.Getenv("HOST")
	user = os.Getenv("USER")
	password = os.Getenv("PASSWORD")
	dbPort = os.Getenv("DB_PORT")
	dbName = os.Getenv("DB_NAME")
}

func StartDB() {
	config := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbName, dbPort)
	dsn := config
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("error connecting to database: ", err)
	}

	fmt.Println("success connect to database")
	db.Debug().AutoMigrate(models.User{}, models.Comment{}, models.Photo{}, models.SocialMedia{})
}

func GetDB() *gorm.DB {
	return db
}
