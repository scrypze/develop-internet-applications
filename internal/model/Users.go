package model

import "github.com/google/uuid"

type Users struct {
	UUID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	// ID          uint      `gorm:"primary_key" json:"id"`

	Login       string    `gorm:"type:varchar(25);unique;not null" json:"login"`
	//Password    string    `gorm:"type:varchar(100);not null" json:"-"`
	//IsModerator bool      `gorm:"type:boolean;default:false" json:"is_moderator"`
	//Name        string    `json:"name"`
	Role        Role      `sql:"type:string;"`
	Pass        string
}
