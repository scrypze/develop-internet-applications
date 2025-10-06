package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) GetSelectedStarsByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		logrus.Error(err)
	}

	selectedStars, err := h.Repository.GetSelectedStarsByID(id)

	if err != nil {
		logrus.Error(err)
	}

	calcByStarID, err := h.Repository.GetCalculateExoplanetsBySelectedStarsID(id)

	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "selected-stars.html", gin.H{
		"selected":        selectedStars,
		"selectedStars":   selectedStars.SelectedStarsItems,
		"selectedStarsID": id,
		"calcByStarID":    calcByStarID,
	})
}

func (h *Handler) AddStarToSelected(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	err = h.Repository.AddStarIntoSelectedStars(id)

	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Redirect(http.StatusSeeOther, "/stars")
}

func (h *Handler) DeleteSelectedStars(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	err = h.Repository.DeleteSelectedStars(id)

	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Redirect(http.StatusSeeOther, "/stars")
}
