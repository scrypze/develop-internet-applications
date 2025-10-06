package handler

import (
	"develop-internet-applications/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct{
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterHandler(router *gin.Engine) {
	router.GET("/stars", h.GetStars)
	router.GET("/star/:id", h.GetStarByID)
	router.GET("/selected-stars/:id", h.GetSelectedStarsByID)
	router.POST("/selected-stars/add-star/:id", h.AddStarToSelected)
	router.POST("/selected-stars/delete-selected-stars/:id", h.DeleteSelectedStars)
}

func (h *Handler) RegisterStatic(router *gin.Engine) {
	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./resources")
}

func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
}
