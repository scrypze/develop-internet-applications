package handler

import (
	"develop-internet-applications/internal/model"
	"fmt"
	"net/http"
	"strconv"

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

// func (h *Handler) Login(ctx *gin.Context) {
// 	var p struct{ Login, Password string }
// 	if err := ctx.BindJSON(&p); err != nil {
// 		h.errorHandler(ctx, http.StatusBadRequest, err)
// 		return
// 	}
// 	u, err := h.service.GetUserByLogin(p.Login)
// 	if err != nil || u.Password != p.Password {
// 		h.errorHandler(ctx, http.StatusUnauthorized, errInvalidCreds())
// 		return
// 	}
// 	ctx.SetCookie("user_id", fmtID(u.ID), 3600, "/", "", false, true)
// 	ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
// }

func (h *Handler) Logout(ctx *gin.Context) {
	ctx.SetCookie("user_id", "", -1, "/", "", false, true)
	ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *Handler) Me(ctx *gin.Context) {
	id, ok := readUserID(ctx)
	if !ok {
		h.errorHandler(ctx, http.StatusUnauthorized, errUnauthorized())
		return
	}
	u, err := h.service.GetUserByID(id)
	if err != nil {
		h.errorHandler(ctx, http.StatusUnauthorized, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"user": u})
}

func fmtID(id uint) string { return fmt.Sprintf("%d", id) }
func readUserID(ctx *gin.Context) (uint, bool) {
	v, err := ctx.Cookie("user_id")
	if err != nil {
		return 0, false
	}
	n, convErr := strconv.ParseUint(v, 10, 64)
	if convErr != nil {
		return 0, false
	}
	return uint(n), true
}
func errRequired() error     { return fmt.Errorf("login and password required") }
func errInvalidCreds() error { return fmt.Errorf("invalid credentials") }
func errUnauthorized() error { return fmt.Errorf("unauthorized") }
