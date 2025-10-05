package repository

import (
	"develop-internet-applications/internal/app/model"

	"gorm.io/gorm/clause"
)

func (r *Repository) GetCalculateExoplanets() (model.CalculateExoplanets, error) {
	var list model.CalculateExoplanets

	err := r.db.Find(&list).Error
	if err != nil {
		return model.CalculateExoplanets{}, err
	}

	return list, nil
}

func (r *Repository) GetCalculateExoplanetsBySelectedStarsID(selectedStarsID int) (map[int]model.CalculateExoplanets, error) {
	var rows []model.CalculateExoplanets
	if err := r.db.Where("selected_stars_id = ?", selectedStarsID).Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make(map[int]model.CalculateExoplanets, len(rows))
	for _, row := range rows {
		result[row.StarID] = row
	}
	return result, nil
}

func (r *Repository) AddStarIntoSelectedStars(id int) error {
	var selectedStarsID int
	creatorID := 1

	err := r.db.Model(&model.SelectedStars{}).
		Where("creator_id = ? AND status = ?", creatorID, "draft").
		Select("id").
		First(&selectedStarsID).Error 

	if err != nil {
		return err
	}

	rec := model.CalculateExoplanets{
		SelectedStarsID: selectedStarsID,
		StarID:          id,
	}

	err = r.db.Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "selected_stars_id"}, {Name: "star_id"}},
			DoNothing: true,
		},
	).Create(&rec).Error
	
	if err != nil {
		return err
	}

	return nil
}
