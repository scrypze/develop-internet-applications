package service

import (
	"develop-internet-applications/internal/model"
	"develop-internet-applications/internal/repository"
	"develop-internet-applications/pkg"
)

type Service struct {
	Star
	SelectedStars
	CalculateExoplanets
	Users
}

type Star interface {
	GetStars() ([]model.Star, error)
	GetStarsByTitle(starTitle string) ([]model.Star, error)
	GetStarByID(id int) (*model.Star, error)
	CreateStar(star *model.Star) (model.Star, error)
	UpdateStar(id int, star *model.Star) error
	DeleteStar(id int) error
	UploadStarImage(id int, file []byte, filename string) error
}

type SelectedStars interface {
	GetSelectedStars() ([]model.SelectedStars, error)
	GetSelectedStarsByID(id int) (model.SelectedStars, error)
	GetSelectedStarsCount() int64
	DeleteSelectedStars(selectedStarsID int) error
	GetCurrentDraftID(creatorID uint) (int, error)
	AddStarIntoSelectedStars(starID int) error
	GetCalculateExoplanetsBySelectedStarsID(selectedStarsID int) (map[int]model.CalculateExoplanets, error)
	GetSelectedStarsFiltered(dateFrom, dateTo, status string) ([]model.SelectedStars, error)
	UpdateSelectedStars(id int, date string, scientist string) error
	FormSelectedStars(id int) error
	RemoveStarFromSelected(starID int) error
	CreateDraftSelectedStars(creatorID int) (model.SelectedStars, error)
	ModerateSelectedStars(id int, moderatorID int, action string) error
	// UpdateCalculateExoplanetsComment(selectedStarsID int, starID int, comment string) error
}

type CalculateExoplanets interface {
	UpdateCalculateExoplanetsComment(selectedStarsID int, starID int, comment string) error
}
type Users interface {
	CreateUser(login, passwordHash string, isModerator bool) (model.Users, error)
	GetUserByLogin(login string) (model.Users, error)
	GetUserByID(id uint) (model.Users, error)
	UpdateUser(id uint, fields map[string]interface{}) error
}

func NewService(repo *repository.Repository, minio *pkg.MinioClient) *Service {
	return &Service{
		Star:                NewStarService(repo.Star, minio),
		SelectedStars:       NewSelectedStarsService(repo.SelectedStars),
		CalculateExoplanets: NewCalculateExoplanetsService(repo.CalculateExoplanets),
		Users:               NewUsersService(repo.Users),
	}
}
