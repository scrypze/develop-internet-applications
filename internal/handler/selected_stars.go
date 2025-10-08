package handler

import (
	"net/http"
	"strconv"
	"time"

	"develop-internet-applications/internal/model"

	"fmt"

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

func (h *Handler) GetSelectedStarsCount(ctx *gin.Context) {
	creatorID := 1

	draftID, err := h.service.GetCurrentDraftID(uint(creatorID))
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	count := h.service.GetSelectedStarsCount()

	ctx.JSON(http.StatusOK, gin.H{
		"selected_stars_id": draftID,
		"count":             count,
	})
}

func (h *Handler) GetSelectedStars(ctx *gin.Context) {
	dateFrom := ctx.Query("date_from")
	dateTo := ctx.Query("date_to")
	status := ctx.Query("status")

	lists, err := h.service.GetSelectedStarsFiltered(dateFrom, dateTo, status)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	type Row struct {
		ID         int     `json:"id"`
		Status     string  `json:"status"`
		FormedAt   string  `json:"formed_at"`
		Creator    string  `json:"creator_login"`
		Moderator  *string `json:"moderator_login"`
		Scientist  string  `json:"scientist"`
		ItemsCount int     `json:"items_count"`
	}

	resp := make([]Row, 0, len(lists))
	for _, l := range lists {
		var modLogin *string
		if l.Moderator.Login != "" {
			ml := l.Moderator.Login
			modLogin = &ml
		}
		row := Row{
			ID:         l.ID,
			Status:     l.Status,
			FormedAt:   l.FormedAt.Format("2006-01-02"),
			Creator:    l.Creator.Login,
			Moderator:  modLogin,
			Scientist:  l.Scientist,
			ItemsCount: len(l.SelectedStarsItems),
		}
		resp = append(resp, row)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"selected-stars": resp,
	})
}

func (h *Handler) RemoveStarFromSelected(ctx *gin.Context) {
	idStr := ctx.Param("id")
	starID, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	if err := h.service.RemoveStarFromSelected(starID); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.Status(http.StatusOK)
}

func (h *Handler) CreateDraftSelectedStars(ctx *gin.Context) {
	creatorID := 1
	draft, err := h.service.CreateDraftSelectedStars(creatorID)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"selected-stars": draft})
}

func (h *Handler) UpdateSelectedStars(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	var payload struct {
		Date      string `json:"date"`
		Scientist string `json:"scientist"`
	}
	if err := ctx.BindJSON(&payload); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	if err := h.service.UpdateSelectedStars(id, payload.Date, payload.Scientist); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	updated, err := h.service.GetSelectedStarsByID(id)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"selected-stars": updated})
}

func (h *Handler) UpdateCalculateExoplanetsComment(ctx *gin.Context) {
	selectedStarsIDStr := ctx.Query("selected_stars_id")
	starIDStr := ctx.Query("star_id")
	selectedStarsID, err := strconv.Atoi(selectedStarsIDStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("bad selected_stars_id"))
		return
	}
	starID, err := strconv.Atoi(starIDStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("bad star_id"))
		return
	}

	var p struct {
		Comment string `json:"comment"`
	}
	if err := ctx.BindJSON(&p); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	if err := h.service.UpdateCalculateExoplanetsComment(selectedStarsID, starID, p.Comment); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.Status(http.StatusOK)
}

func (h *Handler) FormSelectedStars(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	if err := h.service.FormSelectedStars(id); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	updated, err := h.service.GetSelectedStarsByID(id)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"selected-stars": updated})
}

func (h *Handler) ModerateSelectedStars(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	var payload struct {
		Action      string `json:"action"`
		ModeratorID int    `json:"moderator_id"`
	}
	if err := ctx.BindJSON(&payload); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	if payload.ModeratorID == 0 || (payload.Action != "complete" && payload.Action != "decline") {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("invalid payload"))
		return
	}

	if err := h.service.ModerateSelectedStars(id, payload.ModeratorID, payload.Action); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	updated, err := h.service.GetSelectedStarsByID(id)
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

	items := make([]StarWithCalc, 0, len(updated.SelectedStarsItems))
	for _, s := range updated.SelectedStarsItems {
		calc := calcByStarID[s.ID]
		items = append(items, StarWithCalc{
			Star:                    s,
			ProbableNumberOfPlanets: calc.ProbableNumberOfPlanets,
			HabitableZone:           calc.HabitableZone,
		})
	}

	base := SelectedStarsBase{
		ID:          updated.ID,
		Status:      updated.Status,
		CreatedAt:   updated.CreatedAt,
		FormedAt:    updated.FormedAt,
		CompletedAt: updated.CompletedAt,
		CreatorID:   updated.CreatorID,
		ModeratorID: updated.ModeratorID,
		Date:        updated.Date,
		Scientist:   updated.Scientist,
		Creator:     updated.Creator,
		Moderator:   updated.Moderator,
	}

	resp := SelectedResponse{SelectedStarsBase: base, SelectedStarsItems: items}
	ctx.JSON(http.StatusOK, gin.H{"selected-stars": resp})
}
