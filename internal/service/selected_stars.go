package service

import (
	"develop-internet-applications/internal/model"
	"develop-internet-applications/internal/repository"
	"fmt"
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

func (s *SelectedStarsService) GetSelectedStarsFiltered(dateFrom, dateTo, status string) ([]model.SelectedStars, error) {
	return s.repo.GetSelectedStarsFiltered(dateFrom, dateTo, status)
}

func (s *SelectedStarsService) UpdateSelectedStars(id int, date string, scientist string) error {
	return s.repo.UpdateSelectedStars(id, date, scientist)
}

func (s *SelectedStarsService) FormSelectedStars(id int) error {
	entity, err := s.repo.GetSelectedStarsByID(id)
	if err != nil {
		return err
	}
	if entity.Scientist == "" || entity.Date.IsZero() {
		return fmt.Errorf("required fields are empty")
	}
	if len(entity.SelectedStarsItems) == 0 {
		return fmt.Errorf("no stars in the selected list")
	}
	return s.repo.FormSelectedStars(id)
}
