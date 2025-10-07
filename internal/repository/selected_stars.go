package repository

import (
	"develop-internet-applications/internal/model"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type SelectedStarsPostgres struct {
	db *gorm.DB
}

func NewSelectedStarsPostgres(db *gorm.DB) *SelectedStarsPostgres {
	return &SelectedStarsPostgres{db: db}
}

func (r *SelectedStarsPostgres) GetSelectedStars() ([]model.SelectedStars, error) {
	var lists []model.SelectedStars

	err := r.db.Preload("Creator").Preload("Moderator").Preload("SelectedStarsItems").
		Where("status = ?", "draft").
		Find(&lists).Error

	if err != nil {
		return nil, err
	}

	if len(lists) == 0 {
		return nil, fmt.Errorf("список выбранных звёзд пуст")
	}

	return lists, nil
}

func (r *SelectedStarsPostgres) GetSelectedStarsByID(id int) (model.SelectedStars, error) {
	var list model.SelectedStars

	err := r.db.Preload("Creator").Preload("Moderator").Preload("SelectedStarsItems").
		First(&list, id).Error

	if err != nil {
		return model.SelectedStars{}, err
	}

	return list, nil
}

func (r *SelectedStarsPostgres) GetSelectedStarsFiltered(dateFrom, dateTo string, status string) ([]model.SelectedStars, error) {
	var lists []model.SelectedStars

	q := r.db.Preload("Creator").Preload("Moderator").Preload("SelectedStarsItems").Model(&model.SelectedStars{})

	q = q.Where("status NOT IN (?)", []string{"draft", "is_delete"})

	if status != "" {
		q = q.Where("status = ?", status)
	}
	if dateFrom != "" {
		q = q.Where("formed_at >= ?", dateFrom)
	}
	if dateTo != "" {
		q = q.Where("formed_at <= ?", dateTo)
	}

	if err := q.Order("id DESC").Find(&lists).Error; err != nil {
		return nil, err
	}
	return lists, nil
}

func (r *SelectedStarsPostgres) GetSelectedStarsCount() int64 {
	var selectedStarsID int
	var count int64
	creatorID := 1

	err := r.db.Model(&model.SelectedStars{}).
		Where("creator_id = ? AND status = ?", creatorID, "draft").
		Order("id DESC").Limit(1).
		Pluck("id", &selectedStarsID).Error

	if err != nil {
		return 0
	}

	if selectedStarsID == 0 {
		return 0
	}

	err = r.db.Model(&model.CalculateExoplanets{}).
		Where("selected_stars_id = ?", selectedStarsID).
		Count(&count).Error

	if err != nil {
		return 0
	}

	return count
}

func (r *SelectedStarsPostgres) DeleteSelectedStars(selectedStarsID int) error {
	err := r.db.Model(&model.SelectedStars{}).
		Where("id = ?", selectedStarsID).
		Updates(map[string]interface{}{
			"status":    "is_delete",
			"formed_at": time.Now(),
		}).Error

	if err != nil {
		return fmt.Errorf("ошибка удаления списка выбранных звёзд с id %d: %w", selectedStarsID, err)
	}

	return nil
}

func (r *SelectedStarsPostgres) GetCurrentDraftID(creatorID uint) (int, error) {
	var id int

	err := r.db.Model(&model.SelectedStars{}).
		Where("creator_id = ? AND status = ?", creatorID, "draft").
		Order("id DESC").
		Select("id").
		Pluck("id", &id).Error

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *SelectedStarsPostgres) UpdateSelectedStars(id int, date string, scientist string) error {
	updates := map[string]interface{}{}
	if date != "" {
		updates["date"] = date
	}
	if scientist != "" {
		updates["scientist"] = scientist
	}
	if len(updates) == 0 {
		return nil
	}
	return r.db.Model(&model.SelectedStars{}).Where("id = ?", id).Updates(updates).Error
}

func (r *SelectedStarsPostgres) FormSelectedStars(id int) error {
	return r.db.Model(&model.SelectedStars{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":    "formed",
		"formed_at": time.Now(),
	}).Error
}

func (r *SelectedStarsPostgres) CreateDraftSelectedStars(creatorID int) (model.SelectedStars, error) {
	draft := model.SelectedStars{
		Status:    "draft",
		CreatorID: creatorID,
		Date:      time.Now(),
	}
	if err := r.db.Create(&draft).Error; err != nil {
		return model.SelectedStars{}, err
	}
	return draft, nil
}

func (r *SelectedStarsPostgres) ModerateSelectedStars(id int, moderatorID int, action string) error {
	status := ""
	switch action {
	case "complete":
		status = "completed"
	case "decline":
		status = "declined"
	default:
		return fmt.Errorf("unknown action: %s", action)
	}

	var moderator model.Users
	if err := r.db.First(&moderator, moderatorID).Error; err != nil {
		return fmt.Errorf("moderator not found")
	}

	var items []model.CalculateExoplanets
	if err := r.db.Where("selected_stars_id = ?", id).Find(&items).Error; err != nil {
		return err
	}
	for _, it := range items {
		var star model.Star
		if err := r.db.First(&star, it.StarID).Error; err != nil {
			continue
		}
		rin, rout, nPlanets, calcErr := estimateHZAndPlanets(star)
		if calcErr != nil {
			continue
		}
		hz := fmt.Sprintf("%.2f-%.2f a.e.", rin, rout)
		_ = r.db.Model(&model.CalculateExoplanets{}).
			Where("selected_stars_id = ? AND star_id = ?", id, it.StarID).
			Updates(map[string]interface{}{
				"habitable_zone":             hz,
				"probable_number_of_planets": float32(math.Round(nPlanets)),
			}).Error
	}

	return r.db.Model(&model.SelectedStars{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":       status,
			"completed_at": time.Now(),
			"moderator_id": moderatorID,
		}).Error
}

func estimateHZAndPlanets(s model.Star) (rinAU, routAU, nPlanets float64, err error) {
	parseLoose := func(str string) (float64, error) {
		re := regexp.MustCompile(`[+-]?[0-9]+(?:[.,][0-9]+)?`)
		m := re.FindString(strings.TrimSpace(str))
		if m == "" {
			return 0, fmt.Errorf("no number")
		}
		m = strings.ReplaceAll(m, ",", ".")
		v, e := strconv.ParseFloat(m, 64)
		return v, e
	}

	var L float64
	if strings.TrimSpace(s.Luminosity) != "" {
		L, err = parseLoose(s.Luminosity)
		if err != nil {
			return
		}
	} else {
		R, e1 := parseLoose(s.Radius)
		T, e2 := parseLoose(s.Temperature)
		if e1 != nil || e2 != nil {
			err = fmt.Errorf("bad radius/temperature")
			return
		}
		L = (R * R) * math.Pow(T/5778.0, 4)
	}

	const SeffIn = 1.107
	const SeffOut = 0.356
	rinAU = math.Sqrt(L / SeffIn)
	routAU = math.Sqrt(L / SeffOut)

	M, e1 := parseLoose(s.Mass)
	FeH, e2 := parseLoose(s.Metallicity)
	if e1 != nil {
		M = 1
	}
	if e2 != nil {
		FeH = 0
	}
	nPlanets = 2.5 * math.Pow(M, 0.8) * math.Pow(10, 0.3*FeH)
	if nPlanets < 0.5 {
		nPlanets = 0.5
	}
	if nPlanets > 10 {
		nPlanets = 10
	}
	return
}
