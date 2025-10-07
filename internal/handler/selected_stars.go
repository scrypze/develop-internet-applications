package handler

import (
	"net/http"
	"strconv"
	"time"

	"develop-internet-applications/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) GetSelectedStarsByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		logrus.Error(err)
	}

	selectedStars, err := h.service.GetSelectedStarsByID(id)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	calcByStarID, err := h.service.GetCalculateExoplanetsBySelectedStarsID(id)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	type StarWithCalc struct {
		model.Star
		ProbableNumberOfPlanets float32 `json:"probable_number_of_planets"`
		HabitableZone           string  `json:"habitable_zone"`
	}

	type SelectedStarsBase struct {
		ID          int         `json:"ID"`
		Status      string      `json:"Status"`
		CreatedAt   time.Time   `json:"CreatedAt"`
		FormedAt    time.Time   `json:"FormedAt"`
		CompletedAt time.Time   `json:"CompletedAt"`
		CreatorID   int         `json:"CreatorID"`
		ModeratorID *int        `json:"ModeratorID"`
		Date        time.Time   `json:"Date"`
		Scientist   string      `json:"Scientist"`
		Creator     model.Users `json:"Creator"`
		Moderator   model.Users `json:"Moderator"`
	}

	type SelectedResponse struct {
		SelectedStarsBase
		SelectedStarsItems []StarWithCalc `json:"selected-stars-items"`
	}

	items := make([]StarWithCalc, 0, len(selectedStars.SelectedStarsItems))
	for _, s := range selectedStars.SelectedStarsItems {
		calc := calcByStarID[s.ID]
		items = append(items, StarWithCalc{
			Star:                    s,
			ProbableNumberOfPlanets: calc.ProbableNumberOfPlanets,
			HabitableZone:           calc.HabitableZone,
		})
	}

	base := SelectedStarsBase{
		ID:          selectedStars.ID,
		Status:      selectedStars.Status,
		CreatedAt:   selectedStars.CreatedAt,
		FormedAt:    selectedStars.FormedAt,
		CompletedAt: selectedStars.CompletedAt,
		CreatorID:   selectedStars.CreatorID,
		ModeratorID: selectedStars.ModeratorID,
		Date:        selectedStars.Date,
		Scientist:   selectedStars.Scientist,
		Creator:     selectedStars.Creator,
		Moderator:   selectedStars.Moderator,
	}

	resp := SelectedResponse{
		SelectedStarsBase:  base,
		SelectedStarsItems: items,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"selected-stars": resp,
	})
}

func (h *Handler) AddStarToSelected(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	err = h.service.AddStarIntoSelectedStars(id)

	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Status(http.StatusCreated)
}

func (h *Handler) DeleteSelectedStars(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	err = h.service.DeleteSelectedStars(id)

	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Status(http.StatusOK)
}
