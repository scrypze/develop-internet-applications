package main

import (
	"develop-internet-applications/internal/app/dsn"
	"develop-internet-applications/internal/app/model"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	_ = godotenv.Load()
	db, err := gorm.Open(postgres.Open(dsn.FromEnv()), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(
		&model.Users{},
		&model.SelectedStars{},
		&model.Star{},
		&model.CalculateExoplanets{},
	)

	if err != nil {
		panic("cant migrate db")
	}
}
