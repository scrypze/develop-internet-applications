package model

import "github.com/google/uuid"

type Users struct {
	UUID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Login string `gorm:"type:varchar(25);unique;not null" json:"login"`
	Role Role `sql:"type:string;"`
	Pass string
}
