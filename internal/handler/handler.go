package handler

import (
	"develop-internet-applications/internal/model"
	"develop-internet-applications/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterHandler(router *gin.Engine) {
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.URL("http://localhost:8080/swagger/doc.json"),
		ginSwagger.DefaultModelsExpandDepth(-1)))

	api := router.Group("/api")

	auth := api.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/register-astronomer", h.RegisterAstronomer)
		auth.POST("/login", h.Login)
	}

	protectedAuth := api.Group("/auth")
	{
		protectedAuth.Use(h.WithAuthCheck()).GET("/me", h.Me)
		protectedAuth.Use(h.WithAuthCheck()).POST("/logout", h.Logout)
	}

	stars := api.Group("/stars")
	{
		stars.GET("", h.GetStars)
		stars.GET("/:id", h.GetStarByID)
		// stars.POST("", h.CreateStar)
		// stars.PUT("/:id", h.UpdateStar)
		// stars.DELETE("/:id", h.DeleteStar)
		// stars.POST("/:id/image", h.UploadStarImage)
	}

	protectedStars := api.Group("/stars")
	{
		protectedStars.Use(h.WithAuthCheck(model.Astronomer)).POST("", h.CreateStar)
		protectedStars.Use(h.WithAuthCheck(model.Astronomer)).PUT("/:id", h.UpdateStar)
		protectedStars.Use(h.WithAuthCheck(model.Astronomer)).DELETE("/:id", h.DeleteStar)
		protectedStars.Use(h.WithAuthCheck(model.Astronomer)).POST("/:id/image", h.UploadStarImage)
	}

	// selectedStars := api.Group("/selected-stars")
	// {
	// 	selectedStars.POST("", h.CreateDraftSelectedStars)
	// 	selectedStars.GET("/:id", h.GetSelectedStarsByID)
	// 	selectedStars.DELETE("/:id", h.DeleteSelectedStars)
	// 	selectedStars.POST("/add-star/:id", h.AddStarToSelected)
	// 	selectedStars.POST("/delete-selected-stars/:id", h.DeleteSelectedStars)
	// 	selectedStars.GET("", h.GetSelectedStars)
	// 	selectedStars.GET("/count", h.GetSelectedStarsCount)
	// 	selectedStars.PUT("/:id", h.UpdateSelectedStars)
	// 	selectedStars.PUT("/:id/form", h.FormSelectedStars)
	// 	selectedStars.PUT("/:id/moderate", h.ModerateSelectedStars)
	// }

	protectedSelectedStars := api.Group("/selected-stars")
	{
		protectedSelectedStars.Use(h.WithAuthCheck(model.Client, model.Astronomer)).POST("", h.CreateDraftSelectedStars)
		protectedSelectedStars.Use(h.WithAuthCheck(model.Client, model.Astronomer)).GET("/:id", h.GetSelectedStarsByID)
		protectedSelectedStars.Use(h.WithAuthCheck(model.Client, model.Astronomer)).DELETE("/:id", h.DeleteSelectedStars)
		protectedSelectedStars.Use(h.WithAuthCheck(model.Client, model.Astronomer)).POST("/add-star/:id", h.AddStarToSelected)
		protectedSelectedStars.Use(h.WithAuthCheck(model.Client, model.Astronomer)).POST("/delete-selected-stars/:id", h.DeleteSelectedStars)
		protectedSelectedStars.Use(h.WithAuthCheck(model.Astronomer)).GET("", h.GetSelectedStars)
		protectedSelectedStars.Use(h.WithAuthCheck(model.Client, model.Astronomer)).GET("/count", h.GetSelectedStarsCount)
		protectedSelectedStars.Use(h.WithAuthCheck(model.Client, model.Astronomer)).PUT("/:id", h.UpdateSelectedStars)
		protectedSelectedStars.Use(h.WithAuthCheck(model.Client, model.Astronomer)).PUT("/:id/form", h.FormSelectedStars)
		protectedSelectedStars.Use(h.WithAuthCheck(model.Astronomer)).PUT("/:id/moderate", h.ModerateSelectedStars)
	}

	// calculateExoplanets := api.Group("/calculate-exoplanets")
	// {
	// 	calculateExoplanets.PUT("/updateStarComment", h.UpdateCalculateExoplanetsComment)
	// 	calculateExoplanets.DELETE("/remove-star/:id", h.RemoveStarFromSelected)
	// }

	protectedCalculateExoplanets := api.Group("calculate-exoplanets")
	{
		protectedCalculateExoplanets.Use(h.WithAuthCheck(model.Client, model.Astronomer)).PUT("/updateStarComment", h.UpdateCalculateExoplanetsComment)
		protectedCalculateExoplanets.Use(h.WithAuthCheck(model.Client, model.Astronomer)).DELETE("/remove-star/:id", h.RemoveStarFromSelected)
	}
}

func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
}
