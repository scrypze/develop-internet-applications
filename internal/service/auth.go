package service

import (
	"develop-internet-applications/internal/model"
	"develop-internet-applications/internal/repository"
	"develop-internet-applications/pkg/config"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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

	err = bcrypt.CompareHashAndPassword([]byte(user.Pass), []byte(password))
	if err != nil {
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

func (s *AuthService) ValidateToken(tokenStr string) (uuid.UUID, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.JWT.SecretKey), nil
	})
	if err != nil || !token.Valid {
		return uuid.UUID{}, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.UUID{}, fmt.Errorf("invalid claims")
	}

	uuidStr, ok := claims["user_uuid"].(string)
	if !ok {
		return uuid.UUID{}, fmt.Errorf("user_uuid not found")
	}

	return uuid.Parse(uuidStr)
}

func (s *AuthService) GetUserFromToken(tokenStr string) (model.Users, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.JWT.SecretKey), nil
	})
	if err != nil || !token.Valid {
		return model.Users{}, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return model.Users{}, fmt.Errorf("invalid claims")
	}

	uidStr, ok := claims["user_uuid"].(string)
	if !ok {
		return model.Users{}, fmt.Errorf("user_uuid not found in claims")
	}

	uid, err := uuid.Parse(uidStr)
	if err != nil {
		return model.Users{}, fmt.Errorf("invalid UUID")
	}

	user, err := s.repo.GetUserByID(uid)
	if err != nil {
		return model.Users{}, err
	}

	return user, nil
}

func (s *AuthService) Register(login, password string) error {
	_, err := s.repo.GetUserByLogin(login)

	if err == nil {
		return errors.New("user already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	newUUID := uuid.New()
	
	_, err = s.repo.CreateUser(newUUID, login, model.Client, string(hashedPassword))

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}
