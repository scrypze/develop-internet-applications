package handler

import (
	"develop-internet-applications/internal/app/repository"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

func (h *Handler) GetStars(ctx *gin.Context) {
	var stars []repository.Star
	var err error

	searchedStar := ctx.Query("searchedStar")
	if searchedStar == "" {
		stars, err = h.Repository.GetStars()
		if err != nil {
			logrus.Error(err)
		}
	} else {
		stars, err = h.Repository.GetStarsByTitle(searchedStar)
		if err != nil {
			logrus.Error(err)
		}
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
		logrus.Error(err)
	}

	from := ctx.Query("from")
	ctx.HTML(http.StatusOK, "star.html", gin.H{
		"star": star,
		"from": from,
	})

}

func (h *Handler) GetSelectedStarsByID(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
	}

	list, err := h.Repository.GetSelectedStarsByID(id)
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "cart.html", gin.H{
		"cart": list.SelectedStarsItems,
	})
}
