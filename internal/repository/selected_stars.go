package repository

import (
	"develop-internet-applications/internal/model"
	"fmt"

	"gorm.io/gorm"
)

type SelectedStarsPostgres struct {
	db *gorm.DB
}

func NewSelectedStarsPostgres(db *gorm.DB) *SelectedStarsPostgres {
	return &SelectedStarsPostgres{db: db}
}

func (r *SelectedStarsPostgres) GetSelectedStars() ([]model.SelectedStars, error) {
	var lists []model.SelectedStars

	err := r.db.Preload("Creator").Preload("Moderator").Preload("SelectedStarsItems").
		Where("status = ?", "draft").
		Find(&lists).Error

	if err != nil {
		return nil, err
	}

	if len(lists) == 0 {
		return nil, fmt.Errorf("список выбранных звёзд пуст")
	}

	return lists, nil
}

func (r *SelectedStarsPostgres) GetSelectedStarsByID(id int) (model.SelectedStars, error) {
	var list model.SelectedStars

	err := r.db.Preload("Creator").Preload("Moderator").Preload("SelectedStarsItems").
		Where("status = ?", "draft").
		First(&list, id).Error

	if err != nil {
		return model.SelectedStars{}, err
	}

	return list, nil
}

func (r *SelectedStarsPostgres) GetSelectedStarsFiltered(dateFrom, dateTo string, status string) ([]model.SelectedStars, error) {
	var lists []model.SelectedStars

	q := r.db.Preload("Creator").Preload("Moderator").Preload("SelectedStarsItems").Model(&model.SelectedStars{})

	q = q.Where("status NOT IN (?)", []string{"draft", "is_delete"})

	if status != "" {
		q = q.Where("status = ?", status)
	}
	if dateFrom != "" {
		q = q.Where("formed_at >= ?", dateFrom)
	}
	if dateTo != "" {
		q = q.Where("formed_at <= ?", dateTo)
	}

	if err := q.Order("id DESC").Find(&lists).Error; err != nil {
		return nil, err
	}
	return lists, nil
}

func (r *SelectedStarsPostgres) GetSelectedStarsCount() int64 {
	var selectedStarsID int
	var count int64
	creatorID := 1

	err := r.db.Model(&model.SelectedStars{}).
		Where("creator_id = ? AND status = ?", creatorID, "draft").
		Order("id DESC").Limit(1).
		Pluck("id", &selectedStarsID).Error

	if err != nil {
		return 0
	}

	if selectedStarsID == 0 {
		return 0
	}

	err = r.db.Model(&model.CalculateExoplanets{}).
		Where("selected_stars_id = ?", selectedStarsID).
		Count(&count).Error

	if err != nil {
		return 0
	}

	return count
}

func (r *SelectedStarsPostgres) DeleteSelectedStars(selectedStarsID int) error {
	err := r.db.Model(&model.SelectedStars{}).
		Where("id = ?", selectedStarsID).
		Update("status", "is_delete").Error

	if err != nil {
		return fmt.Errorf("ошибка удаления списка выбранных звёзд с id %d: %w", selectedStarsID, err)
	}

	return nil
}

func (r *SelectedStarsPostgres) GetCurrentDraftID(creatorID uint) (int, error) {
	var id int

	err := r.db.Model(&model.SelectedStars{}).
		Where("creator_id = ? AND status = ?", creatorID, "draft").
		Order("id DESC").
		Select("id").
		Pluck("id", &id).Error

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *SelectedStarsPostgres) UpdateSelectedStars(id int, date string, scientist string) error {
	updates := map[string]interface{}{}
	if date != "" {
		updates["date"] = date
	}
	if scientist != "" {
		updates["scientist"] = scientist
	}
	if len(updates) == 0 {
		return nil
	}
	return r.db.Model(&model.SelectedStars{}).Where("id = ?", id).Updates(updates).Error
}
