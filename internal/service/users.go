package service

import (
	"develop-internet-applications/internal/model"
	"develop-internet-applications/internal/repository"
)

type UsersService struct {
	repo repository.Users
}

func NewUsersService(r repository.Users) *UsersService { return &UsersService{repo: r} }

func (s *UsersService) CreateUser(login, passwordHash string, isModerator bool) (model.Users, error) {
	return s.repo.CreateUser(login, passwordHash, isModerator)
}

func (s *UsersService) GetUserByLogin(login string) (model.Users, error) {
	return s.repo.GetUserByLogin(login)
}

func (s *UsersService) GetUserByID(id uint) (model.Users, error) {
	return s.repo.GetUserByID(id)
}

func (s *UsersService) UpdateUser(id uint, fields map[string]interface{}) error {
	return s.repo.UpdateUser(id, fields)
}
