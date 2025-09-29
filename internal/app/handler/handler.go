package handler

import (
	"develop-internet-applications/internal/app/repository"
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

func (h *Handler) RegisterHandler(router *gin.Engine) {
	router.GET("/stars", h.GetStars)
	router.GET("/star/:id", h.GetStar)
	router.GET("/selected-stars/:id", h.GetSelectedStarsByID)
}

func (h *Handler) RegisterStatic(router *gin.Engine) {
	router.LoadHTMLGlob("templates/*")
	router.Static("/styles", "./styles")
}

func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
}
