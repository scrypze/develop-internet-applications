package handler

import (
	"net/http"
	"strconv"
	"time"

	"develop-internet-applications/internal/app/ds"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) GetStars(ctx *gin.Context) {
	var stars []ds.Star
	var err error

	searchedStar := ctx.Query("searchedStar")
	if searchedStar == "" {
		stars, err = h.Repository.GetStars()
	} else {
		stars, err = h.Repository.GetStarsByTitle(searchedStar)
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		logrus.Error(err)
		return
	}

	selectedStarsID := 1

	selectedStars, _ := h.Repository.GetSelectedStarsByID(selectedStarsID)
	selectedStarsSize := len(selectedStars.SelectedStarsItems)

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"time":              time.Now().Format("15:04:05"),
		"stars":             stars,
		"searchedStar":      searchedStar,
		"selectedStarsSize": selectedStarsSize,
		"selectedStarsID":   selectedStarsID,
	})
}

func (h *Handler) GetStar(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
	}

	star, err := h.Repository.GetStar(id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		logrus.Error(err)
		return
	}

	from := ctx.Query("from")
	selectedStarsID := ctx.Query("selectedStarsID")
	ctx.HTML(http.StatusOK, "star.html", gin.H{
		"star":            star,
		"from":            from,
		"selectedStarsID": selectedStarsID,
	})
}
