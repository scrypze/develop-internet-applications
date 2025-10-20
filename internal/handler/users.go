package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// func (h *Handler) Logout(ctx *gin.Context) {
// 	ctx.SetCookie("user_id", "", -1, "/", "", false, true)
// 	ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
// }

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
