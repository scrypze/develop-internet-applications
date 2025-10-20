package handler

import (
	"develop-internet-applications/internal/model"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Login godoc
// @Summary Авторизация пользователя
// @Description Аутентификация пользователя с получением JWT токена
// @Tags Auth
// @Accept json
// @Produce json
// @Param input body model.LoginReq true "Данные для входа"
// @Success 200 {object} model.LoginResp
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /auth/login [post]
func (h *Handler) Login(ctx *gin.Context) {
	var req model.LoginReq

	err := ctx.BindJSON(&req)

	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	token, err := h.service.AuthenticateUser(req.Login, req.Password)
	if err != nil {
		h.errorHandler(ctx, http.StatusUnauthorized, err)
		return
	}

	if token == "" {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("token is nil"))
		return
	}

	ctx.JSON(http.StatusOK, model.LoginResp{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   3600,
	})

}

// Register godoc
// @Summary Регистрация нового пользователя
// @Description Создание нового пользователя с ролью Client
// @Tags Auth
// @Accept json
// @Produce json
// @Param input body model.RegisterReq true "Данные для регистрации"
// @Success 201 {object} model.RegisterResp
// @Failure 400 {object} map[string]interface{}
// @Router /auth/register [post]
func (h *Handler) Register(ctx *gin.Context) {
	var req model.RegisterReq

	err := ctx.BindJSON(&req)

	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	if req.Login == "" || req.Password == "" {
		h.errorHandler(ctx, http.StatusBadRequest, errRequired())
		return
	}

	err = h.service.Register(req.Login, req.Password)

	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusCreated, &model.RegisterResp{
		Message: "user registered successfully",
	})
}

// RegisterAstronomer godoc
// @Summary Регистрация астронома
// @Description Создание нового пользователя с ролью Astronomer
// @Tags Auth
// @Accept json
// @Produce json
// @Param input body model.RegisterReq true "Данные для регистрации"
// @Success 201 {object} model.RegisterResp
// @Failure 400 {object} map[string]interface{}
// @Router /auth/register-astronomer [post]
func (h *Handler) RegisterAstronomer(ctx *gin.Context) {
	var req model.RegisterReq

	err := ctx.BindJSON(&req)

	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	if req.Login == "" || req.Password == "" {
		h.errorHandler(ctx, http.StatusBadRequest, errRequired())
		return
	}

	err = h.service.RegisterAstronomer(req.Login, req.Password)

	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusCreated, &model.RegisterResp{
		Message: "user registered successfully",
	})
}

// Logout godoc
// @Summary Выход пользователя
// @Description Завершение сессии пользователя и инвалидация токена
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /auth/logout [post]
func (h *Handler) Logout(ctx *gin.Context) {
	token, exists := ctx.Get("jwt_token")
	if !exists {
		h.errorHandler(ctx, http.StatusUnauthorized, fmt.Errorf("token not found in context"))
		return
	}

	tokenStr, ok := token.(string)
	if !ok {
		h.errorHandler(ctx, http.StatusInternalServerError, fmt.Errorf("invalid token type"))
		return
	}

	err := h.service.Logout(tokenStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "logged out successfully",
	})
}
