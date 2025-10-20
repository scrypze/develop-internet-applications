package config

import (
	"fmt"
	"os"

	"develop-internet-applications/internal/model"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConfDB() string {
	host := os.Getenv("DB_HOST")
	if host == "" {
		return ""
	}
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	dbname := os.Getenv("DB_NAME")

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, dbname)
}

func MigrateDB() {
	_ = godotenv.Load()
	db, err := gorm.Open(postgres.Open(ConfDB()), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)

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
