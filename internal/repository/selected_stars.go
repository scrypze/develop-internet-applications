package repository

import (
	"develop-internet-applications/internal/model"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (r *SelectedStarsPostgres) GetSelectedStarsFiltered(dateFrom, dateTo string, status string, creatorID uuid.UUID, role model.Role) ([]model.SelectedStars, error) {
	var lists []model.SelectedStars

	q := r.db.Preload("Creator").Preload("Moderator").Preload("SelectedStarsItems").Model(&model.SelectedStars{})

	q = q.Where("status NOT IN (?)", []string{"draft", "is_delete"})

	fmt.Printf("GetSelectedStarsFiltered - role: %d, Astronomer: %d, creatorID: %s\n", role, model.Astronomer, creatorID)
	if role != model.Astronomer {
		fmt.Printf("Filtering by creator_id: %s\n", creatorID)
		q = q.Where("creator_id = ?", creatorID)
	} else {
		fmt.Printf("Astronomer - returning all applications\n")
	}

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
	fmt.Printf("GetSelectedStarsFiltered - found %d applications\n", len(lists))
	return lists, nil
}

func (r *SelectedStarsPostgres) GetSelectedStarsCount(creatorID uuid.UUID) int64 {
	var selectedStarsID int
	var count int64

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

func (r *SelectedStarsPostgres) GetCurrentDraftIDByUUID(creatorUUID uuid.UUID) (int, error) {
	var id int

	err := r.db.Model(&model.SelectedStars{}).
		Where("creator_id = ? AND status = ?", creatorUUID, "draft").
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

func (r *SelectedStarsPostgres) CreateDraftSelectedStars(creatorID uuid.UUID) (model.SelectedStars, error) {
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

func (r *SelectedStarsPostgres) ModerateSelectedStars(id int, moderatorID uuid.UUID, action string) error {
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

	return r.db.Model(&model.SelectedStars{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":       status,
			"completed_at": time.Now(),
			"moderator_id": moderatorID,
		}).Error
}

func (r *SelectedStarsPostgres) AddStarIntoSelectedStars(starID int, creatorID uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var draft model.SelectedStars

		err := tx.Where("creator_id = ? AND status = ?", creatorID, "draft").
			Order("id DESC").Limit(1).
			Find(&draft).Error

		if err != nil {
			return err
		}

		if draft.ID == 0 {
			draft = model.SelectedStars{
				Status:    "draft",
				CreatorID: creatorID,
				Date:      time.Now(),
			}

			err := tx.Create(&draft).Error

			if err != nil {
				return err
			}
		}

		rec := model.CalculateExoplanets{
			SelectedStarsID: draft.ID,
			StarID:          starID,
		}

		return tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "selected_stars_id"}, {Name: "star_id"}},
			DoNothing: true,
		}).Create(&rec).Error
	})
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
