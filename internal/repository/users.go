package repository

import (
	"develop-internet-applications/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UsersPostgres struct{ db *gorm.DB }

func NewUsersPostgres(db *gorm.DB) *UsersPostgres { return &UsersPostgres{db: db} }

func (r *UsersPostgres) CreateUser(uuid uuid.UUID, login string, role model.Role, passwordHash string) (model.Users, error) {
	u := model.Users{UUID: uuid, Login: login, Role: role, Pass: passwordHash}
	if err := r.db.Create(&u).Error; err != nil {
		return model.Users{}, err
	}
	return u, nil
}

func (r *UsersPostgres) GetUserByLogin(login string) (model.Users, error) {
	var u model.Users
	if err := r.db.Where("login = ?", login).First(&u).Error; err != nil {
		return model.Users{}, err
	}
	return u, nil
}

func (r *UsersPostgres) GetUserByID(uuid uuid.UUID) (model.Users, error) {
	var u model.Users
	if err := r.db.First(&u, uuid).Error; err != nil {
		return model.Users{}, err
	}
	return u, nil
}

func (r *UsersPostgres) UpdateUser(uuid uuid.UUID, fields map[string]interface{}) error {
	if len(fields) == 0 {
		return nil
	}
	return r.db.Model(&model.Users{}).Where("uuid = ?", uuid).Updates(fields).Error
}
