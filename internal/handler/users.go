package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Me godoc
// @Summary Получение информации о текущем пользователе
// @Description Возвращает данные авторизованного пользователя
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /auth/me [get]
func (h *Handler) Me(ctx *gin.Context) {
	uuidVal, exists := ctx.Get("user_uuid")
	if !exists {
		h.errorHandler(ctx, http.StatusUnauthorized, errUnauthorized())
		return
	}

	uuid := uuidVal.(uuid.UUID)

	user, err := h.service.GetUserByID(uuid)
	if err != nil {
		h.errorHandler(ctx, http.StatusUnauthorized, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"UUID": user.UUID, "login": user.Login, "role": user.Role})
}

func errRequired() error     { return fmt.Errorf("login and password required") }
func errUnauthorized() error { return fmt.Errorf("unauthorized") }
