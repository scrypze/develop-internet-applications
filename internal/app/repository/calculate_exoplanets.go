package repository

import (
	"develop-internet-applications/internal/app/model"
	"time"

	"gorm.io/gorm"
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

func (r *Repository) AddStarIntoSelectedStars(starID int) error {
	creatorID := 1

	return r.db.Transaction(func(tx *gorm.DB) error {
		var draft model.SelectedStars

		err := tx.Where("creator_id = ? AND status = ?", creatorID, "draft").
			Order("id DESC").Limit(1).
			Find(&draft).Error

		if err != nil {
			return err
		}

		if draft.ID == 0 {
			draft = model.SelectedStars{
				Status:    "draft",
				CreatorID: creatorID,
				Date:      time.Now(),
			}

			err := tx.Create(&draft).Error

			if err != nil {
				return err
			}
		}

		rec := model.CalculateExoplanets{
			SelectedStarsID: draft.ID,
			StarID:          starID,
		}

		return tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "selected_stars_id"}, {Name: "star_id"}},
			DoNothing: true,
		}).Create(&rec).Error
	})
}
