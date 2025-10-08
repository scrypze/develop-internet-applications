package handler

import (
	"develop-internet-applications/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterHandler(router *gin.Engine) {
	api := router.Group("/api")

	auth := api.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
		auth.POST("/logout", h.Logout)
		auth.GET("/me", h.Me)
	}

	stars := api.Group("/stars")
	{
		stars.GET("", h.GetStars)
		stars.GET("/:id", h.GetStarByID)
		stars.POST("", h.CreateStar)
		stars.PUT("/:id", h.UpdateStar)
		stars.DELETE("/:id", h.DeleteStar)
		stars.POST("/:id/image", h.UploadStarImage)
	}

	selectedStars := api.Group("/selected-stars")
	{
		selectedStars.POST("", h.CreateDraftSelectedStars)
		selectedStars.GET("/:id", h.GetSelectedStarsByID)
		selectedStars.DELETE("/:id", h.DeleteSelectedStars)
		selectedStars.POST("/add-star/:id", h.AddStarToSelected)
		selectedStars.DELETE("/remove-star/:id", h.RemoveStarFromSelected)
		selectedStars.PUT("/updateStarComment", h.UpdateCalculateExoplanetsComment)
		selectedStars.POST("/delete-selected-stars/:id", h.DeleteSelectedStars)
		selectedStars.GET("", h.GetSelectedStars)
		selectedStars.GET("/count", h.GetSelectedStarsCount)
		selectedStars.PUT("/:id", h.UpdateSelectedStars)
		selectedStars.PUT("/:id/form", h.FormSelectedStars)
		selectedStars.PUT("/:id/moderate", h.ModerateSelectedStars)
	}
}

func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
}
