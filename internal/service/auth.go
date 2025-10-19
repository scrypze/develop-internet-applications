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

	if user.Password != password {
		return "", errors.New("invalid credentials")
	}

	return s.GenerateJWTToken(user.ID)
}

func (s *AuthService) GenerateJWTToken(userID uint) (string, error) {

	claims := &model.JWTClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "exocalc-app",
		},
		UserUUID: uuid.New(),       
		Scopes:   []string{"user"}, 
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secretKey := "your-secret-key" 
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
