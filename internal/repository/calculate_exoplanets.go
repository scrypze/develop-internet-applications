package repository

import (
	"develop-internet-applications/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
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

func (r *CalculateExoplanetsPostgres) RemoveStarFromSelected(starID int, creatorID uuid.UUID) error {
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
