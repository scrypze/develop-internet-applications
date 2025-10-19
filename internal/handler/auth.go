package handler

import (
	"develop-internet-applications/internal/model"
	"fmt"
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
