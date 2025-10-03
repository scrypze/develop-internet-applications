package repository

import (
	"develop-internet-applications/internal/app/model"
)

func (r *Repository) GetStars() ([]model.Star, error) {
	var stars []model.Star

	err := r.db.Find(&stars).Error

	if err != nil {
		return nil, err
	}

	return stars, nil
}

func (r *Repository) GetStar(id int) (model.Star, error) {
	var star model.Star

	err := r.db.First(&star, id).Error

	if err != nil {
		return model.Star{}, err
	}

	return star, nil
}

func (r *Repository) GetStarsByTitle(title string) ([]model.Star, error) {
	var stars []model.Star

	err := r.db.Where("title ILIKE ?", "%"+title+"%").Find(&stars).Error

	if err != nil {
		return nil, err
	}

	return stars, nil
}
