package handler

import (
	"net/http"
	"strconv"

	"develop-internet-applications/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) GetStars(ctx *gin.Context) {
	var stars []model.Star
	var err error

	searchedStar := ctx.Query("searchedStar")
	if searchedStar == "" {
		stars, err = h.service.GetStars()
	} else {
		stars, err = h.service.GetStarsByTitle(searchedStar)
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		logrus.Error(err)
		return
	}

	creatorID := 1
	selectedStarsID, err := h.service.GetCurrentDraftID(uint(creatorID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"stars":              stars,
		"selectedStarsCount": h.service.GetSelectedStarsCount(),
		"searchedStar":       searchedStar,
		"selectedStarsID":    selectedStarsID,
	})
}

func (h *Handler) GetStarByID(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
	}

	star, err := h.service.GetStarByID(id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		logrus.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"star": star,
	})
}
