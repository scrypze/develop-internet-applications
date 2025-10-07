package repository

import (
	"develop-internet-applications/internal/model"

	"gorm.io/gorm"
)

type Repository struct {
	Star
	SelectedStars
}

type Star interface {
	GetStars() ([]model.Star, error)
	GetStarsByTitle(starTitle string) ([]model.Star, error)
	GetStarByID(id int) (*model.Star, error)
	CreateStar(star *model.Star) (model.Star, error)
	UpdateStar(id int, star *model.Star) error
	DeleteStar(id int) error
	UpdateStarImage(id int, imagePath string) error
}

type SelectedStars interface {
	GetSelectedStars() ([]model.SelectedStars, error)
	GetSelectedStarsByID(id int) (model.SelectedStars, error)
	GetSelectedStarsCount() int64
	DeleteSelectedStars(selectedStarsID int) error
	GetCurrentDraftID(creatorID uint) (int, error)
	AddStarIntoSelectedStars(starID int) error
	GetCalculateExoplanetsBySelectedStarsID(selectedStarsID int) (map[int]model.CalculateExoplanets, error)
	GetSelectedStarsFiltered(dateFrom, dateTo string, status string) ([]model.SelectedStars, error)
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		Star:          NewStarPostgres(db),
		SelectedStars: NewSelectedStarsPostgres(db),
	}
}
