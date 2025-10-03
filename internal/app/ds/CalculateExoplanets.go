package ds

type CalculateExoplanets struct {
	SelectionStarsID        int     `gorm:"primaryKey;column:selection_stars_id"`
	StarID                  int     `gorm:"primaryKey;column:star_id"`
	Comment                 string  `gorm:"type:text"`
	ProbableNumberOfPlanets float32 `gorm:"type:real"`
	HabitableZone           string  `gorm:"type:varchar(64)"`

	SelectedStars SelectedStars `gorm:"foreignKey:SelectionStarsID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Star          Star          `gorm:"foreignKey:StarID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
