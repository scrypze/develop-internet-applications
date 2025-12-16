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

func (s *CalculateExoplanetsService) UpdateCalculateExoplanetsResult(selectedStarsID int, starID int, habitableZone string, probableNumberOfPlanets float32) error {
	return s.repo.UpdateCalculateExoplanetsResult(selectedStarsID, starID, habitableZone, probableNumberOfPlanets)
}

func (s *CalculateExoplanetsService) GetCalculatedCount(selectedStarsID int) (int64, error) {
	return s.repo.GetCalculatedCount(selectedStarsID)
}
