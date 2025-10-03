package ds

type CalculateExoplanets struct {
	SelectionStarsID int    `gorm:"not null;uniqueIndex:idx_sel_star"`
	StarID           int    `gorm:"not null;uniqueIndex:idx_sel_star"`
	Comment                 string  `gorm:"type:text"`
	
	ProbableNumberOfPlanets float32 `gorm:"type:real"`
	HabitableZone           string  `gorm:"type:varchar(64)"`

	SelectedStars SelectedStars `gorm:"foreignKey:SelectionStarsID"`
	Star          Star          `gorm:"foreignKey:StarID"`
}
