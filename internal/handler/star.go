package handler

import (
	"io"
	"net/http"
	"strconv"

	"develop-internet-applications/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// GetStars godoc
// @Summary Получение списка звезд
// @Description Возвращает список всех звезд или отфильтрованный по названию
// @Tags Stars
// @Accept json
// @Produce json
// @Param searchedStar query string false "Поиск по названию звезды"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /stars [get]
func (h *Handler) GetStars(ctx *gin.Context) {
	var stars []model.Star
	var err error

	searchedStar := ctx.Query("searchedStar")
	if searchedStar == "" {
		stars, err = h.service.GetStars()
	} else {
		stars, err = h.service.GetStarsByTitle(searchedStar)
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		logrus.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"stars": stars,
	})
}

// GetStarByID godoc
// @Summary Получение звезды по ID
// @Description Возвращает данные звезды по указанному ID
// @Tags Stars
// @Accept json
// @Produce json
// @Param id path int true "ID звезды"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /stars/{id} [get]
func (h *Handler) GetStarByID(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
	}

	star, err := h.service.GetStarByID(id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		logrus.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"star": star,
	})
}

// CreateStar godoc
// @Summary Создание новой звезды
// @Description Создает новую звезду (требуется роль Astronomer)
// @Tags Stars
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body model.Star true "Данные звезды"
// @Success 201 {object} model.Star
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /stars [post]
func (h *Handler) CreateStar(c *gin.Context) {
	var star model.Star
	if err := c.BindJSON(&star); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid input data")
		return
	}

	createdStar, err := h.service.Star.CreateStar(&star)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, createdStar)
}

// UpdateStar godoc
// @Summary Обновление звезды
// @Description Обновляет данные существующей звезды (требуется роль Astronomer)
// @Tags Stars
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID звезды"
// @Param input body model.Star true "Обновленные данные звезды"
// @Success 200 {object} model.Star
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /stars/{id} [put]
func (h *Handler) UpdateStar(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid star ID")
		return
	}

	var payload model.Star
	if err := c.BindJSON(&payload); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid input data")
		return
	}

	if err := h.service.UpdateStar(id, &payload); err != nil {
		if err.Error() == "star not found" {
			newErrorResponse(c, http.StatusNotFound, err.Error())
		} else {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	updated, _ := h.service.GetStarByID(id)
	c.JSON(http.StatusOK, updated)
}

// DeleteStar godoc
// @Summary Удаление звезды
// @Description Удаляет звезду по ID (требуется роль Astronomer)
// @Tags Stars
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID звезды"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /stars/{id} [delete]
func (h *Handler) DeleteStar(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid star ID")
		return
	}

	if err := h.service.DeleteStar(id); err != nil {
		if err.Error() == "star not found" {
			newErrorResponse(c, http.StatusNotFound, err.Error())
		} else {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	c.Status(http.StatusOK)
}

// UploadStarImage godoc
// @Summary Загрузка изображения звезды
// @Description Загружает изображение для звезды (требуется роль Astronomer)
// @Tags Stars
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID звезды"
// @Param image formData file true "Файл изображения"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /stars/{id}/image [post]
func (h *Handler) UploadStarImage(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid star ID")
		return
	}

	file, header, err := c.Request.FormFile("image")
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "No image file provided")
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Failed to read file")
		return
	}

	if err := h.service.UploadStarImage(id, fileBytes, header.Filename); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Image uploaded successfully"})
}
