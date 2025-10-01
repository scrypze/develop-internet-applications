package ds

type CalculateExoplanets struct {
	SelectionStarsID int    `gorm:"primaryKey;column:selection_stars_id"` // FK -> selected_stars.id
	StarID           int    `gorm:"primaryKey;column:star_id"`            // FK -> stars.id
	Comment          string `gorm:"type:text"`
}

func (CalculateExoplanets) TableName() string { return "calculate_exoplanets" }
