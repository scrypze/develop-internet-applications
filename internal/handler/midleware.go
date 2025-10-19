package handler

import (
	"develop-internet-applications/internal/model"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (h *Handler) WithAuthCheck(ctx *gin.Context) {
	jwtStr := ctx.GetHeader("Authorization")

	if !strings.HasPrefix(jwtStr, model.JwtPrefix) {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	jwtStr = jwtStr[len(model.JwtPrefix):]

	uuid, err := h.service.ValidateToken(jwtStr)
	if err != nil {
		ctx.AbortWithStatus(http.StatusForbidden)
		log.Println(err)
		return
	}

	ctx.Set("user_uuid", uuid)

	ctx.Next()
}
