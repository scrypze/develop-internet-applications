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
	router.Use(h.CORSMiddleware())

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
		protectedAuth.GET("/me", h.WithAuthCheck(), h.Me)
		protectedAuth.POST("/logout", h.WithAuthCheck(), h.Logout)
		protectedAuth.PUT("/update-login", h.WithAuthCheck(), h.UpdateLogin)
		protectedAuth.PUT("/update-password", h.WithAuthCheck(), h.UpdatePassword)
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
		protectedStars.POST("", h.WithAuthCheck(model.Astronomer), h.CreateStar)
		protectedStars.PUT("/:id", h.WithAuthCheck(model.Astronomer), h.UpdateStar)
		protectedStars.DELETE("/:id", h.WithAuthCheck(model.Astronomer), h.DeleteStar)
		protectedStars.POST("/:id/image", h.WithAuthCheck(model.Astronomer), h.UploadStarImage)
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
		protectedSelectedStars.POST("", h.WithAuthCheck(model.Client, model.Astronomer), h.CreateDraftSelectedStars)
		protectedSelectedStars.GET("/:id", h.WithAuthCheck(model.Client, model.Astronomer), h.GetSelectedStarsByID)
		protectedSelectedStars.DELETE("/:id", h.WithAuthCheck(model.Client, model.Astronomer), h.DeleteSelectedStars)
		protectedSelectedStars.POST("/add-star/:id", h.WithAuthCheck(model.Client, model.Astronomer), h.AddStarToSelected)
		protectedSelectedStars.POST("/delete-selected-stars/:id", h.WithAuthCheck(model.Client, model.Astronomer), h.DeleteSelectedStars)
		protectedSelectedStars.GET("", h.WithAuthCheck(model.Client, model.Astronomer), h.GetSelectedStars)
		protectedSelectedStars.GET("/count", h.GetSelectedStarsCount)
		protectedSelectedStars.PUT("/:id", h.WithAuthCheck(model.Client, model.Astronomer), h.UpdateSelectedStars)
		protectedSelectedStars.PUT("/:id/form", h.WithAuthCheck(model.Client, model.Astronomer), h.FormSelectedStars)
		protectedSelectedStars.PUT("/:id/moderate", h.WithAuthCheck(model.Astronomer), h.ModerateSelectedStars)
	}

	// calculateExoplanets := api.Group("/calculate-exoplanets")
	// {
	// 	calculateExoplanets.PUT("/updateStarComment", h.UpdateCalculateExoplanetsComment)
	// 	calculateExoplanets.DELETE("/remove-star/:id", h.RemoveStarFromSelected)
	// }

	protectedCalculateExoplanets := api.Group("calculate-exoplanets")
	{
		protectedCalculateExoplanets.PUT("/updateStarComment", h.WithAuthCheck(model.Client, model.Astronomer), h.UpdateCalculateExoplanetsComment)
		protectedCalculateExoplanets.DELETE("/remove-star/:id", h.WithAuthCheck(model.Client, model.Astronomer), h.RemoveStarFromSelected)
	}
}

func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
}
