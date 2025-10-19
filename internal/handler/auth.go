package handler

import (
	"develop-internet-applications/internal/model"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

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

func (h *Handler) WithAuthCheck(ctx *gin.Context) {
	jwtStr := ctx.GetHeader("Authorization")

	if !strings.HasPrefix(jwtStr, model.JwtPrefix) {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	jwtStr = jwtStr[len(model.JwtPrefix):]

	_, err := h.service.ValidateToken(jwtStr)
	if err != nil {
		ctx.AbortWithStatus(http.StatusForbidden)
		log.Println(err)
		return
	}
}
