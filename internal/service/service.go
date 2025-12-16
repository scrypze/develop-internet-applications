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
	GetSelectedStarsCount(creatorID uuid.UUID) int64
	DeleteSelectedStars(selectedStarsID int) error
	GetCurrentDraftID(creatorID uint) (int, error)
	GetCurrentDraftIDByUUID(creatorUUID uuid.UUID) (int, error)
	AddStarIntoSelectedStars(starID int, creatorID uuid.UUID) error
	GetCalculateExoplanetsBySelectedStarsID(selectedStarsID int) (map[int]model.CalculateExoplanets, error)
	GetSelectedStarsFiltered(dateFrom, dateTo, status string, creatorID uuid.UUID, role model.Role) ([]model.SelectedStars, error)
	UpdateSelectedStars(id int, date string, scientist string) error
	FormSelectedStars(id int) error
	CreateDraftSelectedStars(creatorID uuid.UUID) (model.SelectedStars, error)
	ModerateSelectedStars(id int, moderatorID uuid.UUID, action string) error
}

type CalculateExoplanets interface {
	UpdateCalculateExoplanetsComment(selectedStarsID int, starID int, comment string) error
	RemoveStarFromSelected(starID int, creatorID uuid.UUID) error
	UpdateCalculateExoplanetsResult(selectedStarsID int, starID int, habitableZone string, probableNumberOfPlanets float32) error
	GetCalculatedCount(selectedStarsID int) (int64, error)
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
	Logout(tokenStr string) error
	RegisterAstronomer(login, password string) error
	HashPassword(password string) (string, error)
}

func NewService(repo *repository.Repository, minio *pkg.MinioClient, redis *pkg.RedisClient, cfg *config.Config) *Service {
	return &Service{
		Star:                NewStarService(repo.Star, minio),
		SelectedStars:       NewSelectedStarsService(repo.SelectedStars, repo.Star, cfg),
		CalculateExoplanets: NewCalculateExoplanetsService(repo.CalculateExoplanets),
		Users:               NewUsersService(repo.Users),
		Auth:                NewAuthService(repo.Users, redis, cfg),
	}
}
