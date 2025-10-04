package repository

import (
	"develop-internet-applications/internal/app/model"
	"fmt"

	"github.com/sirupsen/logrus"
)

func (r *Repository) GetSelectedStars() ([]model.SelectedStars, error) {
	var lists []model.SelectedStars
	if err := r.db.Preload("SelectedStarsItems").Find(&lists).Error; err != nil {
		return nil, err
	}
	if len(lists) == 0 {
		return nil, fmt.Errorf("список выбранных звёзд пуст")
	}
	return lists, nil
}

func (r *Repository) GetSelectedStarsByID(id int) (model.SelectedStars, error) {
	var list model.SelectedStars
	if err := r.db.Preload("SelectedStarsItems").First(&list, id).Error; err != nil {
		return model.SelectedStars{}, err
	}
	return list, nil
}

func (r *Repository) GetSelectedStarsCount() int64 {
	var selectedStarsID int
	var count int64

	creatorID := 1

	err := r.db.Model(&model.SelectedStars{}).
		Where("creator_id = ? AND status = ?", creatorID, "draft").
		Select("id").
		First(&selectedStarsID).Error

	if err != nil {
		return 0
	}

	err = r.db.Model(&model.CalculateExoplanets{}).
		Where("selected_stars_id = ?", selectedStarsID).
		Count(&count).Error

	if err != nil {
		logrus.Error("error counting records in calculate_exoplanets:", err)
		return 0
	}

	return count
}
