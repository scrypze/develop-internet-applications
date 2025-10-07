package service

import (
	"develop-internet-applications/internal/model"
	"develop-internet-applications/internal/repository"
)

type SelectedStarsService struct {
	repo repository.SelectedStars
}

func NewSelectedStarsService(repo repository.SelectedStars) *SelectedStarsService {
	return &SelectedStarsService{repo: repo}
}

func (s *SelectedStarsService) GetSelectedStars() ([]model.SelectedStars, error) {
	return s.repo.GetSelectedStars()
}

func (s *SelectedStarsService) GetSelectedStarsByID(id int) (model.SelectedStars, error) {
	return s.repo.GetSelectedStarsByID(id)
}

func (s *SelectedStarsService) GetSelectedStarsCount() int64 {
	return s.repo.GetSelectedStarsCount()
}

func (s *SelectedStarsService) DeleteSelectedStars(selectedStarsID int) error {
	return s.repo.DeleteSelectedStars(selectedStarsID)
}

func (s *SelectedStarsService) GetCurrentDraftID(creatorID uint) (int, error) {
	return s.repo.GetCurrentDraftID(creatorID)
}

func (s *SelectedStarsService) AddStarIntoSelectedStars(starID int) error {
	return s.repo.AddStarIntoSelectedStars(starID)
}

func (s *SelectedStarsService) GetCalculateExoplanetsBySelectedStarsID(selectedStarsID int) (map[int]model.CalculateExoplanets, error) {
	return s.repo.GetCalculateExoplanetsBySelectedStarsID(selectedStarsID)
}
