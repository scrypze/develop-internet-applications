package handler

import (
	"net/http"
	"strconv"
	"time"

	"develop-internet-applications/internal/model"

	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// GetSelectedStarsByID godoc
// @Summary Получение заявки по ID
// @Description Возвращает заявку со всеми расчетами экзопланет
// @Tags SelectedStars
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID заявки"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /selected-stars/{id} [get]
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
		ID             int        `json:"ID"`
		Status         string     `json:"Status"`
		CreatedAt      time.Time  `json:"CreatedAt"`
		FormedAt       time.Time  `json:"FormedAt"`
		CompletedAt    time.Time  `json:"CompletedAt"`
		CreatorID      uuid.UUID  `json:"CreatorID"`
		ModeratorID    *uuid.UUID `json:"ModeratorID"`
		Date           time.Time  `json:"Date"`
		Scientist      string     `json:"Scientist"`
		//CreatorLogin   string     `json:"creator_login"`
		//ModeratorLogin *string    `json:"moderator_login"`
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

	// var modLogin *string
	// if selectedStars.Moderator.Login != "" {
	// 	ml := selectedStars.Moderator.Login
	// 	modLogin = &ml
	// }
	base := SelectedStarsBase{
		ID:             selectedStars.ID,
		Status:         selectedStars.Status,
		CreatedAt:      selectedStars.CreatedAt,
		FormedAt:       selectedStars.FormedAt,
		CompletedAt:    selectedStars.CompletedAt,
		CreatorID:      selectedStars.CreatorID,
		ModeratorID:    selectedStars.ModeratorID,
		Date:           selectedStars.Date,
		Scientist:      selectedStars.Scientist,
		//CreatorLogin:   selectedStars.Creator.Login,
		//ModeratorLogin: modLogin,
	}

	resp := SelectedResponse{
		SelectedStarsBase:  base,
		SelectedStarsItems: items,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"selected-stars": resp,
	})
}

// AddStarToSelected godoc
// @Summary Добавление звезды в заявку
// @Description Добавляет звезду по ID в текущую заявку
// @Tags SelectedStars
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID звезды"
// @Success 201
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /selected-stars/add-star/{id} [post]
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

// DeleteSelectedStars godoc
// @Summary Удаление заявки
// @Description Удаляет заявку по ID
// @Tags SelectedStars
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID заявки"
// @Success 200
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /selected-stars/{id} [delete]
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

// GetSelectedStarsCount godoc
// @Summary Получение количества звезд в текущей заявке
// @Description Возвращает ID текущей черновой заявки и количество звезд в ней
// @Tags SelectedStars
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /selected-stars/count [get]
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

// GetSelectedStars godoc
// @Summary Получение списка заявок
// @Description Возвращает отфильтрованный список заявок
// @Tags SelectedStars
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param date_from query string false "Дата начала (YYYY-MM-DD)"
// @Param date_to query string false "Дата окончания (YYYY-MM-DD)"
// @Param status query string false "Статус заявки (draft, formed, completed, declined)"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /selected-stars [get]
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

// RemoveStarFromSelected godoc
// @Summary Удаление звезды из заявки
// @Description Удаляет звезду из текущей заявки
// @Tags CalculateExoplanets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID звезды"
// @Success 200
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /calculate-exoplanets/remove-star/{id} [delete]
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

// CreateDraftSelectedStars godoc
// @Summary Создание черновика заявки
// @Description Создает новую черновую заявку
// @Tags SelectedStars
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 201 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /selected-stars [post]
func (h *Handler) CreateDraftSelectedStars(ctx *gin.Context) {
	creatorID, _ := uuid.Parse("b57f6d40-23a8-4e8c-9a14-1d2d2fa68a6b")

	draft, err := h.service.CreateDraftSelectedStars(creatorID)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"selected-stars": draft})
}

// UpdateSelectedStars godoc
// @Summary Обновление заявки
// @Description Обновляет дату и ученого для заявки
// @Tags SelectedStars
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID заявки"
// @Param input body map[string]interface{} true "Дата и ученый"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /selected-stars/{id} [put]
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

// FormSelectedStars godoc
// @Summary Формирование заявки
// @Description Переводит заявку из черновика в статус "formed"
// @Tags SelectedStars
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID заявки"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /selected-stars/{id}/form [put]
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

// ModerateSelectedStars godoc
// @Summary Модерация заявки
// @Description Модератор подтверждает или отклоняет заявку (требуется роль Astronomer)
// @Tags SelectedStars
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID заявки"
// @Param input body map[string]interface{} true "Действие (complete/decline) и ID модератора"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /selected-stars/{id}/moderate [put]
func (h *Handler) ModerateSelectedStars(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	var payload struct {
		Action      string `json:"action"`
		ModeratorID uuid.UUID    `json:"moderator_id"`
	}
	if err := ctx.BindJSON(&payload); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	if payload.ModeratorID == uuid.Nil || (payload.Action != "complete" && payload.Action != "decline") {
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
		ID             int        `json:"ID"`
		Status         string     `json:"Status"`
		CreatedAt      time.Time  `json:"CreatedAt"`
		FormedAt       time.Time  `json:"FormedAt"`
		CompletedAt    time.Time  `json:"CompletedAt"`
		CreatorID      uuid.UUID  `json:"CreatorID"`
		ModeratorID    *uuid.UUID `json:"ModeratorID"`
		Date           time.Time  `json:"Date"`
		Scientist      string     `json:"Scientist"`
		CreatorLogin   string     `json:"creator_login"`
		ModeratorLogin *string    `json:"moderator_login"`
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

	var modLogin *string
	if updated.Moderator.Login != "" {
		ml := updated.Moderator.Login
		modLogin = &ml
	}
	base := SelectedStarsBase{
		ID:             updated.ID,
		Status:         updated.Status,
		CreatedAt:      updated.CreatedAt,
		FormedAt:       updated.FormedAt,
		CompletedAt:    updated.CompletedAt,
		CreatorID:      updated.CreatorID,
		ModeratorID:    updated.ModeratorID,
		Date:           updated.Date,
		Scientist:      updated.Scientist,
		CreatorLogin:   updated.Creator.Login,
		ModeratorLogin: modLogin,
	}

	resp := SelectedResponse{SelectedStarsBase: base, SelectedStarsItems: items}
	ctx.JSON(http.StatusOK, gin.H{"selected-stars": resp})
}
