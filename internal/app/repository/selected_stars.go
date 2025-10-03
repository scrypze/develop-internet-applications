package repository

import (
	"develop-internet-applications/internal/app/model"
	"fmt"
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
