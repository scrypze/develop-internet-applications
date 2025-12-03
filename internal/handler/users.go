package handler

import (
	"develop-internet-applications/internal/model"
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

	roleStr := "client"
	if user.Role == model.Astronomer {
		roleStr = "astronomer"
	} else if user.Role == model.Guest {
		roleStr = "guest"
	}
	ctx.JSON(http.StatusOK, gin.H{"UUID": user.UUID, "login": user.Login, "role": roleStr})
}

// UpdateLogin godoc
// @Summary Обновление логина пользователя
// @Description Обновляет логин текущего авторизованного пользователя
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body map[string]string true "Новый логин"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /auth/update-login [put]
func (h *Handler) UpdateLogin(ctx *gin.Context) {
	uuidVal, exists := ctx.Get("user_uuid")
	if !exists {
		h.errorHandler(ctx, http.StatusUnauthorized, errUnauthorized())
		return
	}

	userUUID := uuidVal.(uuid.UUID)

	var payload struct {
		Login string `json:"login" binding:"required"`
	}
	if err := ctx.BindJSON(&payload); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("invalid request body"))
		return
	}

	if payload.Login == "" {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("login cannot be empty"))
		return
	}

	// Проверяем, не занят ли логин другим пользователем
	existingUser, err := h.service.GetUserByLogin(payload.Login)
	if err == nil && existingUser.UUID != userUUID {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("login already taken"))
		return
	}

	fields := map[string]interface{}{
		"login": payload.Login,
	}

	if err := h.service.UpdateUser(userUUID, fields); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Login updated successfully"})
}

// UpdatePassword godoc
// @Summary Обновление пароля пользователя
// @Description Обновляет пароль текущего авторизованного пользователя
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body map[string]string true "Старый и новый пароль"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /auth/update-password [put]
func (h *Handler) UpdatePassword(ctx *gin.Context) {
	uuidVal, exists := ctx.Get("user_uuid")
	if !exists {
		h.errorHandler(ctx, http.StatusUnauthorized, errUnauthorized())
		return
	}

	userUUID := uuidVal.(uuid.UUID)

	var payload struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}
	if err := ctx.BindJSON(&payload); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("invalid request body"))
		return
	}

	if payload.NewPassword == "" || len(payload.NewPassword) < 6 {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("new password must be at least 6 characters"))
		return
	}

	// Получаем пользователя
	user, err := h.service.GetUserByID(userUUID)
	if err != nil {
		h.errorHandler(ctx, http.StatusUnauthorized, err)
		return
	}

	// Проверяем старый пароль
	_, err = h.service.AuthenticateUser(user.Login, payload.OldPassword)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("invalid old password"))
		return
	}

	// Хэшируем новый пароль
	hashedPassword, err := h.service.HashPassword(payload.NewPassword)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, fmt.Errorf("failed to hash password"))
		return
	}

	fields := map[string]interface{}{
		"pass": hashedPassword,
	}

	if err := h.service.UpdateUser(userUUID, fields); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}

func errRequired() error     { return fmt.Errorf("login and password required") }
func errUnauthorized() error { return fmt.Errorf("unauthorized") }
