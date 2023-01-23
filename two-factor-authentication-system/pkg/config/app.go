package config

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func Connect() {
	dsn := "host=localhost user=postgres password=root dbname=tfas port=5433 sslmode=disable TimeZone=Asia/Dhaka"
	d, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(d)
	}

	db = d
}

func GetDB() *gorm.DB {
	return db
}
