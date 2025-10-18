package handler

import (
	"develop-internet-applications/internal/model"
	"net/http"

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

	ctx.JSON(http.StatusOK, model.LoginResp{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   3600,
	})
}
