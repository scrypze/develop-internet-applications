package service

import (
	"develop-internet-applications/internal/model"
	"develop-internet-applications/internal/repository"
	"develop-internet-applications/pkg/config"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type AuthService struct {
	repo   repository.Users
	config *config.Config
}

func NewAuthService(repo repository.Users, config *config.Config) *AuthService {
	return &AuthService{
		repo:   repo,
		config: config,
	}
}

func (s *AuthService) AuthenticateUser(login, password string) (string, error) {
	user, err := s.repo.GetUserByLogin(login)
	if err != nil {
		return "", err
	}

	if user.Pass != password {
		return "", errors.New("invalid credentials")
	}

	return s.GenerateJWTToken(user.UUID)
}

func (s *AuthService) GenerateJWTToken(userUUID uuid.UUID) (string, error) {

	claims := &model.JWTClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "exocalc-app",
		},
		UserUUID: userUUID,
		Scopes:   []string{"user"},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(s.config.JWT.SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *AuthService) ValidateToken(tokenString string) (*model.JWTClaims, error) {
	claims := &model.JWTClaims{}

	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.JWT.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	return claims, nil
}
