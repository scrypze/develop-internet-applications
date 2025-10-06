package model

import "time"

type SelectedStars struct {
	ID          int       `gorm:"primaryKey;column:id;autoIncrement"`
	Status      string    `gorm:"type:varchar(64)"`
	CreatedAt   time.Time `gorm:"type:date"`
	FormedAt    time.Time `gorm:"type:date"`
	CompletedAt time.Time `gorm:"type:date"`
	CreatorID   int       `gorm:"column:creator_id;not null"`
	ModeratorID *int      `gorm:"column:moderator_id;null"`
	Date        time.Time `gorm:"type:date"`
	Scientist   string    `gorm:"type:varchar(64)"`

	Creator            Users  `gorm:"foreignKey:CreatorID"`
	Moderator          Users  `gorm:"foreignKey:ModeratorID"`
	SelectedStarsItems []Star `gorm:"many2many:calculate_exoplanets;joinForeignKey:SelectedStarsID;joinReferences:StarID"`
}
