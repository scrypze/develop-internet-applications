package repository

import (
	"develop-internet-applications/internal/app/ds"
)

func (r *Repository) GetStars() ([]ds.Star, error) {
	var stars []ds.Star

	err := r.db.Find(&stars).Error

	if err != nil {
		return nil, err
	}

	return stars, nil
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
