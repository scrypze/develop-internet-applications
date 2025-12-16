package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UpdateCalculateExoplanetsComment godoc
// @Summary Обновление комментария для расчетов экзопланет
// @Description Обновляет комментарий для конкретной звезды в заявке
// @Tags CalculateExoplanets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param selected_stars_id query int true "ID заявки"
// @Param star_id query int true "ID звезды"
// @Param input body map[string]interface{} true "Комментарий"
// @Success 200
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /calculate-exoplanets/updateStarComment [put]
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

// UpdateCalculateExoplanetsResult godoc
// @Summary Обновление результата расчета экзопланет
// @Description Обновляет результат расчета для конкретной звезды в заявке (вызывается асинхронным сервисом)
// @Tags CalculateExoplanets
// @Accept json
// @Produce json
// @Param input body map[string]interface{} true "Данные результата"
// @Success 200
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /calculate-exoplanets/update-result [post]
func (h *Handler) UpdateCalculateExoplanetsResult(ctx *gin.Context) {
	token := ctx.GetHeader("X-Service-Token")
	expectedToken := h.config.ComputingServiceToken
	if token != expectedToken {
		h.errorHandler(ctx, http.StatusUnauthorized, fmt.Errorf("invalid service token"))
		return
	}

	var payload struct {
		SelectedStarsID         int     `json:"selected_stars_id"`
		StarID                  int     `json:"star_id"`
		HabitableZone           string  `json:"habitable_zone"`
		ProbableNumberOfPlanets float32 `json:"probable_number_of_planets"`
	}
	if err := ctx.BindJSON(&payload); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	if err := h.service.UpdateCalculateExoplanetsResult(
		payload.SelectedStarsID,
		payload.StarID,
		payload.HabitableZone,
		payload.ProbableNumberOfPlanets,
	); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Status(http.StatusOK)
}
