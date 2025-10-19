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
	Password string `json:"password"`
}

type LoginResp struct {
	ExpiresIn   int    `json:"expires_in"`
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

type JWTClaims struct {
	jwt.StandardClaims           // все что точно необходимо по RFC
	UserUUID           uuid.UUID `json:"user_uuid"` // наши данные - uuid этого пользователя в базе данных
	Scopes             []string  `json:"scopes"`    // список доступов в нашей системе
}
