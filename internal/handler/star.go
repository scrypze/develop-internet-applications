package handler

import (
	"net/http"
	"strconv"

	"develop-internet-applications/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

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
