package database

import (
	"fmt"
	"hello_world/model"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var err error

func Init() {
	DatabaseConnection()
}

func DatabaseConnection() {
	host := "psql"
	port := "5432"
	dbName := "bills"
	dbUser := "postgres"
	password := "pass1234"
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host,
		dbUser,
		password,
		dbName,
		port,
	)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	DB.AutoMigrate(&model.Bill{})
	if err != nil {
		log.Fatal("Error connecting to the database", err)
	}
	fmt.Println("Database connection successful")
}
