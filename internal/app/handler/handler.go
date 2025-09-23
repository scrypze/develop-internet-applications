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

	_, cartCount, _ := h.Repository.GetCart()
	cartCountStr := strconv.Itoa(cartCount)

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"time":      time.Now().Format("15:04:05"),
		"stars":     stars,
		"query":     searchQuery,
		"cartCount": cartCountStr,
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

	ctx.HTML(http.StatusOK, "star.html", gin.H{
		"star": star,
	})

}

func (h *Handler) GetCart(ctx *gin.Context) {
	cart, cartCount, err := h.Repository.GetCart()
	if err != nil {
		logrus.Error(err)
	}
	ctx.HTML(http.StatusOK, "cart.html", gin.H{
		"cart":      cart,
		"cartCount": cartCount,
	})
}
