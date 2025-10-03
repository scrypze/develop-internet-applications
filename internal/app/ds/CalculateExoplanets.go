package ds

type CalculateExoplanets struct {
	SelectedStarsID         int           `gorm:"not null;column:selected_stars_id;uniqueIndex:idx_sel_star"`
	StarID                  int           `gorm:"not null;uniqueIndex:idx_sel_star"`
	Comment                 string        `gorm:"type:text"`
	
	ProbableNumberOfPlanets float32       `gorm:"type:real"`
	HabitableZone           string        `gorm:"type:varchar(64)"`

	SelectedStars           SelectedStars `gorm:"foreignKey:SelectedStarsID"`
	Star                    Star          `gorm:"foreignKey:StarID"`
}
