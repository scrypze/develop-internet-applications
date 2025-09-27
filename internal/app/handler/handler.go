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

	searchQuery := ctx.Query("query")
	if searchQuery == "" {
		stars, err = h.Repository.GetStars()
		if err != nil {
			logrus.Error(err)
		}
	} else {
		stars, err = h.Repository.GetStarsByTitle(searchQuery)
		if err != nil {
			logrus.Error(err)
		}
	}

	starsCartID := 1

	starsCart, _ := h.Repository.GetStarsCartByID(starsCartID)
	starsCartSize := len(starsCart.StarsCartItems)

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"time":      time.Now().Format("15:04:05"),
		"stars":     stars,
		"query":     searchQuery,
		"starsCartSize": starsCartSize,
		"starsCartID": starsCartID,
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

func (h *Handler) GetStarsCartByID(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
	}

	cart, err := h.Repository.GetStarsCartByID(id)
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "cart.html", gin.H{
		"cart":      cart.StarsCartItems,
	})
}
