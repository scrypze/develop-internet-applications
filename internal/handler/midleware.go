package handler

import (
	"develop-internet-applications/internal/model"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (h *Handler) WithAuthCheck(allowedRoles ...model.Role) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		jwtStr := ctx.GetHeader("Authorization")

		if !strings.HasPrefix(jwtStr, model.JwtPrefix) {
			ctx.AbortWithStatus(http.StatusForbidden)
			log.Println("missing or invalid Authorization header")
			return
		}

		jwtStr = jwtStr[len(model.JwtPrefix):]

		ctx.Set("jwt_token", jwtStr)

		if len(allowedRoles) == 0 {
			uuid, err := h.service.ValidateToken(jwtStr)
			if err != nil {
				ctx.AbortWithStatus(http.StatusForbidden)
				log.Printf("token validation failed: %v", err)
				return
			}

			ctx.Set("user_uuid", uuid)
			ctx.Next()
			return
		}

		claims, err := h.service.ValidateTokenWithRole(jwtStr)
		
		if err != nil {
			ctx.AbortWithStatus(http.StatusForbidden)
			log.Printf("token validation failed: %v", err)
			return
		}

		roleAllowed := false
		for _, allowedRole := range allowedRoles {
			if claims.Role == allowedRole {
				roleAllowed = true
				break
			}
		}

		if !roleAllowed {
			ctx.AbortWithStatus(http.StatusForbidden)
			log.Printf("role %d is not allowed, required one of: %v", claims.Role, allowedRoles)
			return
		}

		ctx.Set("user_uuid", claims.UserUUID)
		ctx.Next()
	}
}
