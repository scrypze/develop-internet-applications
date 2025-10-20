package model

import (
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

const (
	JwtPrefix = "Bearer "
)

type LoginReq struct {
	Login    string `json:"login"`
	Password string `json:"pass"`
}

type LoginResp struct {
	ExpiresIn   int    `json:"expires_in"`
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

type RegisterReq struct {
	Login    string `json:"login"`
	Password string `json:"pass"`
}

type RegisterResp struct {
	Message string `json:"message"`
}

type JWTClaims struct {
	jwt.StandardClaims
	UserUUID uuid.UUID `json:"user_uuid"`
	Role     Role      `json:"role"`
}
