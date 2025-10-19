package handler

import (
	"develop-internet-applications/internal/model"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) Register(ctx *gin.Context) {
	var p struct{ Login, Password string }
	if err := ctx.BindJSON(&p); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	if p.Login == "" || p.Password == "" {
		h.errorHandler(ctx, http.StatusBadRequest, errRequired())
		return
	}
	user, err := h.service.CreateUser(uuid.New(), p.Login, model.Client, p.Password)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"user": user})
}

func (h *Handler) Logout(ctx *gin.Context) {
	ctx.SetCookie("user_id", "", -1, "/", "", false, true)
	ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
}

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
