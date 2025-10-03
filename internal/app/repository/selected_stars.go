package repository

import (
	"develop-internet-applications/internal/app/ds"
	"fmt"
)

func (r *Repository) GetSelectedStars() ([]ds.SelectedStars, error) {
	var lists []ds.SelectedStars
	if err := r.db.Preload("SelectedStarsItems").Find(&lists).Error; err != nil {
		return nil, err
	}
	if len(lists) == 0 {
		return nil, fmt.Errorf("список выбранных звёзд пуст")
	}
	return lists, nil
}

func (r *Repository) GetSelectedStarsByID(id int) (ds.SelectedStars, error) {
	var list ds.SelectedStars
	if err := r.db.Preload("SelectedStarsItems").First(&list, id).Error; err != nil {
		return ds.SelectedStars{}, err
	}
	return list, nil
}