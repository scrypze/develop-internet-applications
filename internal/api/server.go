package api

import (
	"log"

	"develop-internet-applications/internal/app/handler"
	"develop-internet-applications/internal/app/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func StartServer() {
	log.Println("Start server")

	repo, err := repository.NewRepository()
	if err != nil {
		logrus.Error("ошибка инициализации репозитория")
	}

	handler := handler.NewHandler(repo)

	r := gin.Default()

	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./resources")

	r.GET("/stars", handler.GetStars)
	r.GET("/star/:id", handler.GetStar)
	r.GET("/cart", handler.GetCart)

	r.Run()
	log.Println("Server down")
}
