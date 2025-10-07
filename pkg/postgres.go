package pkg

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresDB(config string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func SyncSequences(db *gorm.DB) error {
	if err := db.Exec("SELECT setval(pg_get_serial_sequence('stars','id'), COALESCE((SELECT MAX(id) FROM stars), 0))").Error; err != nil {
		return err
	}
	if err := db.Exec("SELECT setval(pg_get_serial_sequence('selected_stars','id'), COALESCE((SELECT MAX(id) FROM selected_stars), 0))").Error; err != nil {
		return err
	}
	return nil
}
