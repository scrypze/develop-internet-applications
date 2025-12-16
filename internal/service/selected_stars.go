package service

import (
	"bytes"
	"develop-internet-applications/internal/model"
	"develop-internet-applications/internal/repository"
	"develop-internet-applications/pkg/config"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type SelectedStarsService struct {
	repo         repository.SelectedStars
	starRepo     repository.Star
	config       *config.Config
}

func NewSelectedStarsService(repo repository.SelectedStars, starRepo repository.Star, cfg *config.Config) *SelectedStarsService {
	return &SelectedStarsService{repo: repo, starRepo: starRepo, config: cfg}
}

func (s *SelectedStarsService) GetSelectedStars() ([]model.SelectedStars, error) {
	return s.repo.GetSelectedStars()
}

func (s *SelectedStarsService) GetSelectedStarsByID(id int) (model.SelectedStars, error) {
	return s.repo.GetSelectedStarsByID(id)
}

func (s *SelectedStarsService) GetSelectedStarsCount(creatorID uuid.UUID) int64 {
	return s.repo.GetSelectedStarsCount(creatorID)
}

func (s *SelectedStarsService) DeleteSelectedStars(selectedStarsID int) error {
	return s.repo.DeleteSelectedStars(selectedStarsID)
}

func (s *SelectedStarsService) GetCurrentDraftID(creatorID uint) (int, error) {
	return s.repo.GetCurrentDraftID(creatorID)
}

func (s *SelectedStarsService) GetCurrentDraftIDByUUID(creatorUUID uuid.UUID) (int, error) {
	return s.repo.GetCurrentDraftIDByUUID(creatorUUID)
}

func (s *SelectedStarsService) AddStarIntoSelectedStars(starID int, creatorID uuid.UUID) error {
	return s.repo.AddStarIntoSelectedStars(starID, creatorID)
}

func (s *SelectedStarsService) GetCalculateExoplanetsBySelectedStarsID(selectedStarsID int) (map[int]model.CalculateExoplanets, error) {
	return s.repo.GetCalculateExoplanetsBySelectedStarsID(selectedStarsID)
}

func (s *SelectedStarsService) GetSelectedStarsFiltered(dateFrom, dateTo, status string, creatorID uuid.UUID, role model.Role) ([]model.SelectedStars, error) {
	return s.repo.GetSelectedStarsFiltered(dateFrom, dateTo, status, creatorID, role)
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

func (s *SelectedStarsService) CreateDraftSelectedStars(creatorID uuid.UUID) (model.SelectedStars, error) {
	return s.repo.CreateDraftSelectedStars(creatorID)
}

func (s *SelectedStarsService) ModerateSelectedStars(id int, moderatorID uuid.UUID, action string) error {
	if _, err := s.repo.GetSelectedStarsByID(id); err != nil {
		return err
	}
	
	// Если действие - завершение заявки, вызываем асинхронный сервис
	if action == "complete" {
		selectedStars, err := s.repo.GetSelectedStarsByID(id)
		if err != nil {
			return err
		}
		
		// Получаем все записи CalculateExoplanets для этой заявки
		calcItems, err := s.repo.GetCalculateExoplanetsBySelectedStarsID(id)
		if err != nil {
			return err
		}
		
		// Для каждой звезды запускаем асинхронный расчет
		for starID := range calcItems {
			go s.callComputingService(id, starID, selectedStars.SelectedStarsItems)
		}
	}
	
	// Обновляем статус заявки (без расчета, расчет будет асинхронным)
	return s.repo.ModerateSelectedStars(id, moderatorID, action)
}

func (s *SelectedStarsService) callComputingService(selectedStarsID int, starID int, stars []model.Star) {
	starData, err := s.starRepo.GetStarByID(starID)
	if err != nil || starData == nil {
		logrus.Errorf("Error getting star data for star_id=%d: %v", starID, err)
		return
	}
	
	payload := map[string]interface{}{
		"selected_stars_id": selectedStarsID,
		"star_id":           starID,
		"star_data": map[string]interface{}{
			"luminosity":  starData.Luminosity,
			"radius":      starData.Radius,
			"temperature": starData.Temperature,
			"mass":        starData.Mass,
			"metallicity": starData.Metallicity,
		},
	}
	
	jsonData, err := json.Marshal(payload)
	if err != nil {
		logrus.Errorf("Error marshaling payload: %v", err)
		return
	}
	
	url := fmt.Sprintf("%s/api/calculate", s.config.ComputingServiceURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		logrus.Errorf("Error creating request: %v", err)
		return
	}
	
	req.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	
	resp, err := client.Do(req)
	if err != nil {
		logrus.Errorf("Error calling computing service: %v", err)
		return
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusAccepted {
		logrus.Errorf("Computing service returned status %d", resp.StatusCode)
		return
	}
	
	logrus.Infof("Successfully called computing service for selected_stars_id=%d, star_id=%d", selectedStarsID, starID)
}
