package service

import (
	"develop-internet-applications/internal/model"
	"develop-internet-applications/internal/repository"
	"develop-internet-applications/pkg"
	"develop-internet-applications/pkg/config"

	"github.com/google/uuid"
)

type Service struct {
	Star
	SelectedStars
	CalculateExoplanets
	Users
	Auth
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
	CreateDraftSelectedStars(creatorID uuid.UUID) (model.SelectedStars, error)
	ModerateSelectedStars(id int, moderatorID int, action string) error
}

type CalculateExoplanets interface {
	UpdateCalculateExoplanetsComment(selectedStarsID int, starID int, comment string) error
	RemoveStarFromSelected(starID int) error
}
type Users interface {
	CreateUser(UUID uuid.UUID, login string, role model.Role, passwordHash string) (model.Users, error)
	GetUserByLogin(login string) (model.Users, error)
	GetUserByID(uuid uuid.UUID) (model.Users, error)
	UpdateUser(uuid uuid.UUID, fields map[string]interface{}) error
}

type Auth interface {
	AuthenticateUser(login, password string) (string, error)
	GenerateJWTToken(userUUID uuid.UUID, role model.Role) (string, error)
	ValidateToken(tokenString string) (uuid.UUID, error)
	ValidateTokenWithRole(tokenString string) (*model.JWTClaims, error)
	Register(login, password string) error
}

func NewService(repo *repository.Repository, minio *pkg.MinioClient, config *config.Config) *Service {
	return &Service{
		Star:                NewStarService(repo.Star, minio),
		SelectedStars:       NewSelectedStarsService(repo.SelectedStars),
		CalculateExoplanets: NewCalculateExoplanetsService(repo.CalculateExoplanets),
		Users:               NewUsersService(repo.Users),
		Auth:                NewAuthService(repo.Users, config),
	}
}
