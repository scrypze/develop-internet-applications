package service

import (
	"develop-internet-applications/internal/model"
	"develop-internet-applications/internal/repository"
	"develop-internet-applications/pkg"
)


type Service struct {
	Star
}

type Star interface {
	GetStars() ([]model.Star, error)
	GetStarsByTitle(starTitle string) ([]model.Star, error)
	GetStarByID(id int) (*model.Star, error)
	CreateStar(star *model.Star) (model.Star, error)
	UpdateStar(id int, star *model.Star) (error)
	DeleteStar(id int) (error)
	UploadMaterialImage(id int, file []byte, filename string) error
}

func NewService(repo *repository.Repository, minio *pkg.MinioClient) *Service {
	return &Service{
		Star: NewStarService(repo.Star, minio),
	}
}
