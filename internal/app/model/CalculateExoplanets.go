package model

type CalculateExoplanets struct {
	SelectedStarsID int    `gorm:"not null;column:selected_stars_id;uniqueIndex:idx_sel_star"`
	StarID          int    `gorm:"not null;uniqueIndex:idx_sel_star"`
	Comment         string `gorm:"type:text"`

	ProbableNumberOfPlanets float32 `gorm:"type:real"`
	HabitableZone           string  `gorm:"type:varchar(64)"`

	SelectedStars SelectedStars `gorm:"foreignKey:SelectedStarsID"`
	Star          Star          `gorm:"foreignKey:StarID"`
}

// r_in(AU) = sqrt( L / 1.107 )
// r_out(AU) = sqrt( L / 0.356 )
// HabitableZone = r_in - r_out a.u.
// ProbableNumberOfPlanets = 2.5 × M^0.8 × 10^(0.3·[Fe/H])
