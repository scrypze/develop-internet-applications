package repository

import (
	"develop-internet-applications/internal/app/ds"
	"fmt"
	"time"
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

func (r *Repository) GetStars() ([]ds.Star, error) {
	var stars []ds.Star

	err := r.db.Find(&stars).Error

	if err != nil {
		return nil, err
	}

	return stars, nil
}

func mustParseDateDDMMYYYY(value string) time.Time {
	t, err := time.Parse("02.01.2006", value)
	if err != nil {
		return time.Time{}
	}
	return t
}

func (r *Repository) GetStar(id int) (ds.Star, error) {
	var star ds.Star

	err := r.db.First(&star, id).Error

	if err != nil {
		return ds.Star{}, err
	}

	return star, nil
}

func (r *Repository) GetStarsByTitle(title string) ([]ds.Star, error) {
	var stars []ds.Star

	err := r.db.Where("title ILIKE ?", "%"+title+"%").Find(&stars).Error

	if err != nil {
		return nil, err
	}

	return stars, nil
}
