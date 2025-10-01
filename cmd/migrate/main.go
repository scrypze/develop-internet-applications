package main

import (
	"develop-internet-applications/internal/app/ds"
	"develop-internet-applications/internal/app/dsn"

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
		&ds.Star{},
		&ds.SelectedStars{},
		&ds.CalculateExoplanets{},
		&ds.Users{},
	)

	if err != nil {
		panic("cant migrate db")
	}
}
