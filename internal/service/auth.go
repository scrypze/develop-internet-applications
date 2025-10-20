package service

import (
	"context"
	"develop-internet-applications/internal/model"
	"develop-internet-applications/internal/repository"
	"develop-internet-applications/pkg"
	"develop-internet-applications/pkg/config"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo        repository.Users
	redisClient *pkg.RedisClient
	config      *config.Config
}

func NewAuthService(repo repository.Users, redisClient *pkg.RedisClient, config *config.Config) *AuthService {
	return &AuthService{
		repo:        repo,
		redisClient: redisClient,
		config:      config,
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

	return s.GenerateJWTToken(user.UUID, user.Role)
}

func (s *AuthService) GenerateJWTToken(userUUID uuid.UUID, role model.Role) (string, error) {
	claims := &model.JWTClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "exocalc-app",
		},
		UserUUID: userUUID,
		Role:     role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(s.config.JWT.SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *AuthService) ValidateToken(tokenStr string) (uuid.UUID, error) {
	ctx := context.Background()
	isBlacklisted, err := s.redisClient.CheckJWTInBlacklist(ctx, tokenStr)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("error checking blacklist: %w", err)
	}
	if isBlacklisted {
		return uuid.UUID{}, fmt.Errorf("token is blacklisted")
	}

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

func (s *AuthService) ValidateTokenWithRole(tokenStr string) (*model.JWTClaims, error) {
	ctx := context.Background()
	isBlacklisted, err := s.redisClient.CheckJWTInBlacklist(ctx, tokenStr)
	if err != nil {
		return nil, fmt.Errorf("error checking blacklist: %w", err)
	}
	if isBlacklisted {
		return nil, fmt.Errorf("token is blacklisted")
	}

	token, err := jwt.ParseWithClaims(tokenStr, &model.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.JWT.SecretKey), nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(*model.JWTClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	return claims, nil
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

func (s *AuthService) RegisterAstronomer(login, password string) error {
	_, err := s.repo.GetUserByLogin(login)

	if err == nil {
		return errors.New("user already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	newUUID := uuid.New()

	_, err = s.repo.CreateUser(newUUID, login, model.Astronomer, string(hashedPassword))

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (s *AuthService) Logout(tokenStr string) error {
	token, err := jwt.ParseWithClaims(tokenStr, &model.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.JWT.SecretKey), nil
	})
	if err != nil {
		return fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*model.JWTClaims)
	if !ok {
		return fmt.Errorf("invalid claims")
	}

	expiresAt := time.Unix(claims.ExpiresAt, 0)
	ttl := time.Until(expiresAt)

	if ttl <= 0 {
		return nil
	}

	ctx := context.Background()
	return s.redisClient.WriteJWTToBlacklist(ctx, tokenStr, ttl)
}
