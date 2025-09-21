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

	searchQuery := ctx.Query("query") // получаем значение из поля поиска
	if searchQuery == "" {            // если поле поиска пусто, то просто получаем из репозитория все записи
		stars, err = h.Repository.GetStars()
		if err != nil {
			logrus.Error(err)
		}
	} else {
		stars, err = h.Repository.GetStarsByTitle(searchQuery) // в ином случае ищем заказ по заголовку
		if err != nil {
			logrus.Error(err)
		}
	}

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"time":  time.Now().Format("15:04:05"),
		"stars": stars,
		"query": searchQuery, // передаем введенный запрос обратно на страницу
		// в ином случае оно будет очищаться при нажатии на кнопку
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
