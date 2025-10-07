package repository

import (
	"database/sql"
	"develop-internet-applications/internal/model"
	"errors"

	"gorm.io/gorm"
)

type StarPostgres struct {
	db *gorm.DB
}

func NewStarPostgres(db *gorm.DB) *StarPostgres {
	return &StarPostgres{db: db}
}

func (r *StarPostgres) GetStars() ([]model.Star, error) {
	var stars []model.Star

	err := r.db.Where("is_deleted = ?", "false").Find(&stars).Error

	if err != nil {
		return nil, err
	}

	return stars, nil
}

func (r *StarPostgres) GetStarByID(id int) (*model.Star, error) {
	query := "SELECT id, title, description, image_path, spectral_type, temperature, radius, mass, luminosity, metallicity, age, distance, is_deleted FROM stars WHERE id = $1"

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
		&star.IsDeleted,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return star, nil
}

func (r *StarPostgres) GetStarsByTitle(starTitle string) ([]model.Star, error) {
	var stars []model.Star

	err := r.db.Where("title ILIKE ? AND is_deleted = ?", "%"+starTitle+"%", "false").Find(&stars).Error

	if err != nil {
		return nil, err
	}

	return stars, nil
}

func (r *StarPostgres) CreateStar(star *model.Star) (model.Star, error) {
	var createdStar model.Star
	createdStar = *star

	err := r.db.Create(&createdStar).Error
	return createdStar, err
}

func (r *StarPostgres) UpdateStar(id int, star *model.Star) error {
	err := r.db.Where("id = ?", id).Updates(star).Error
	return err
}

func (r *StarPostgres) DeleteStar(id int) error {
	return r.db.Model(&model.Star{}).Where("id = ?", id).Update("is_deleted", "true").Error
}

func (r *StarPostgres) UpdateStarImage(id int, imagePath string) error {
	return r.db.Model(&model.Star{}).Where("id = ?", id).Update("image_path", imagePath).Error
}
