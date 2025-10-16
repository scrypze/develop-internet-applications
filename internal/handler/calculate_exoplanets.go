package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

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
