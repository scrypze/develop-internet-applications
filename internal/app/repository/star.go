package repository

import (
	"database/sql"
	"develop-internet-applications/internal/app/model"
	"errors"
)

func (r *Repository) GetStars() ([]model.Star, error) {
	var stars []model.Star

	err := r.db.Find(&stars).Error

	if err != nil {
		return nil, err
	}

	return stars, nil
}

func (r *Repository) GetStarByID(id int) (*model.Star, error) {
	query := "SELECT id, title, description, image_path, spectral_type, temperature, radius, mass, luminosity, metallicity, age, distance FROM stars WHERE id = $1"

	row := r.db.Raw(query, id).Row()

	star := &model.Star{}

	err := row.Scan(
		&star.ID,
		&star.Title,
		&star.Description,
		&star.ImagePath,
		&star.SpectralType,
		&star.Temperature,
		&star.Radius,
		&star.Mass,
		&star.Luminosity,
		&star.Metallicity,
		&star.Age,
		&star.Distance,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
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
