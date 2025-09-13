package handler

import (
	"develop-internet-applications/internal/app/repository"
	"net/http"
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
	stars, err := h.Repository.GetStars()

	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"time": time.Now().Format("15:04:05"),
		"stars": stars,
	})
} 
