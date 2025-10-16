package repository

import (
	"develop-internet-applications/internal/model"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CalculateExoplanetsPostgres struct {
	db *gorm.DB
}

func NewCalculateExoplanetsPostgres(db *gorm.DB) *CalculateExoplanetsPostgres {
	return &CalculateExoplanetsPostgres{db: db}
}

func (r *SelectedStarsPostgres) GetCalculateExoplanets() (model.CalculateExoplanets, error) {
	var list model.CalculateExoplanets

	err := r.db.Find(&list).Error
	if err != nil {
		return model.CalculateExoplanets{}, err
	}

	return list, nil
}

func (r *SelectedStarsPostgres) GetCalculateExoplanetsBySelectedStarsID(selectedStarsID int) (map[int]model.CalculateExoplanets, error) {
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

func (r *SelectedStarsPostgres) AddStarIntoSelectedStars(starID int) error {
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

func (r *SelectedStarsPostgres) RemoveStarFromSelected(starID int) error {
	creatorID := 1
	return r.db.Transaction(func(tx *gorm.DB) error {
		var draft model.SelectedStars
		if err := tx.Where("creator_id = ? AND status = ?", creatorID, "draft").Order("id DESC").Limit(1).Find(&draft).Error; err != nil {
			return err
		}
		if draft.ID == 0 {
			return nil
		}
		return tx.Where("selected_stars_id = ? AND star_id = ?", draft.ID, starID).Delete(&model.CalculateExoplanets{}).Error
	})
}

func (r *CalculateExoplanetsPostgres) UpdateCalculateExoplanetsComment(selectedStarsID int, starID int, comment string) error {
	return r.db.Model(&model.CalculateExoplanets{}).
		Where("selected_stars_id = ? AND star_id = ?", selectedStarsID, starID).
		Update("comment", comment).Error
}
