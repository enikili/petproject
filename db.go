package main

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// переменная, через которую мы будем работать с БД
var DB *gorm.DB

func InitDB() {
	dsn := "host=localhost user=postgres password=12345 dbname=postgres port=5438 sslmode=disable"
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}
}
