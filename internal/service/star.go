package service

import (
	"bytes"
	"develop-internet-applications/internal/model"
	"develop-internet-applications/internal/repository"
	"develop-internet-applications/pkg"
	"fmt"
)

type StarService struct {
	repo        repository.Star
	minioClient *pkg.MinioClient
}

func NewStarService(repo repository.Star, minioClient *pkg.MinioClient) *StarService {
	return &StarService{
		repo:        repo,
		minioClient: minioClient,
	}
}

func (s *StarService) GetStars() ([]model.Star, error) {
	return s.repo.GetStars()
}

func (s *StarService) GetStarByID(id int) (*model.Star, error) {
	star, err := s.repo.GetStarByID(id)
	if err != nil {
		return &model.Star{}, err
	}
	if star == nil || star.ID == 0 {
		return &model.Star{}, fmt.Errorf("star not found")
	}
	return star, nil
}

func (s *StarService) GetStarsByTitle(starTitle string) ([]model.Star, error) {
	return s.repo.GetStarsByTitle(starTitle)
}

func (s *StarService) CreateStar(star *model.Star) (model.Star, error) {
	star.ID = 0
	return s.repo.CreateStar(star)
}

func (s *StarService) UpdateStar(id int, payload *model.Star) error {
	existing, _ := s.repo.GetStarByID(id)
	if existing.ID == 0 {
		return fmt.Errorf("star not found")
	}

	payload.ID = 0
	return s.repo.UpdateStar(id, payload)
}

func (s *StarService) DeleteStar(id int) error {
	star, err := s.repo.GetStarByID(id)

	if star.ID == 0 {
		return fmt.Errorf("star not found: %s", err)
	}

	if star.ImagePath != "" {
		if err := s.minioClient.DeleteImage(id); err != nil {
			fmt.Printf("Warning: failed to delete image from Minio: %v\n", err)
		}
	}

	return s.repo.DeleteStar(id)
}

func (s *StarService) UploadStarImage(id int, file []byte, filename string) error {
	star, err := s.repo.GetStarByID(id)
	if star.ID == 0 {
		return fmt.Errorf("star not found")
	}

	reader := bytes.NewReader(file)
	imageURL, err := s.minioClient.UploadImage(star.Title, reader, int64(len(file)), filename)
	if err != nil {
		return err
	}

	return s.repo.UpdateStarImage(id, imageURL)
}
