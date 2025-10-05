package repository

import (
	"develop-internet-applications/internal/app/model"
	"errors"
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

        err := tx.Where("creator_id = ? AND status = ?", creatorID, "draft").First(&draft).Error

        if errors.Is(err, gorm.ErrRecordNotFound) {
            draft = model.SelectedStars{
                Status:    "draft",
                CreatorID: creatorID,
                CreatedAt:      time.Now(),
            }

            err := tx.Create(&draft).Error

			if err != nil {
                return err
            }

        } else if err != nil {
            return err
        }

        rec := model.CalculateExoplanets{
            SelectedStarsID: draft.ID,
            StarID:          starID,
        }
		
        if err := tx.Clauses(
            clause.OnConflict{
                Columns:   []clause.Column{{Name: "selected_stars_id"}, {Name: "star_id"}},
                DoNothing: true,
            },
        ).Create(&rec).Error; err != nil {
            return err
        }
        return nil
    })
}
