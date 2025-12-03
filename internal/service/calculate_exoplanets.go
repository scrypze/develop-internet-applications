package service

import (
	"develop-internet-applications/internal/repository"

	"github.com/google/uuid"
)

type CalculateExoplanetsService struct {
	repo repository.CalculateExoplanets
}

func NewCalculateExoplanetsService(repo repository.CalculateExoplanets) *CalculateExoplanetsService {
	return &CalculateExoplanetsService{repo: repo}
}

func (s *CalculateExoplanetsService) UpdateCalculateExoplanetsComment(selectedStarsID int, starID int, comment string) error {
	return s.repo.UpdateCalculateExoplanetsComment(selectedStarsID, starID, comment)
}

func (s *CalculateExoplanetsService) RemoveStarFromSelected(starID int, creatorID uuid.UUID) error {
	return s.repo.RemoveStarFromSelected(starID, creatorID)
}
