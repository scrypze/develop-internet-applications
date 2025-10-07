package repository

import (
	"develop-internet-applications/internal/model"

	"gorm.io/gorm"
)

type UsersPostgres struct{ db *gorm.DB }

func NewUsersPostgres(db *gorm.DB) *UsersPostgres { return &UsersPostgres{db: db} }

func (r *UsersPostgres) CreateUser(login, passwordHash string, isModerator bool) (model.Users, error) {
	u := model.Users{Login: login, Password: passwordHash, IsModerator: isModerator}
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

func (r *UsersPostgres) GetUserByID(id uint) (model.Users, error) {
	var u model.Users
	if err := r.db.First(&u, id).Error; err != nil {
		return model.Users{}, err
	}
	return u, nil
}

func (r *UsersPostgres) UpdateUser(id uint, fields map[string]interface{}) error {
	if len(fields) == 0 {
		return nil
	}
	return r.db.Model(&model.Users{}).Where("id = ?", id).Updates(fields).Error
}
