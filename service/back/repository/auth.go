package repository

import (
	"github.com/eavlongs/file_sync/models"
	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) CreateUser(user *models.User) error {
	err := r.db.Create(user).Error

	if err != nil {
		return err
	}

	return r.FindUser(user.Email, user)
}

func (r *AuthRepository) FindUser(email string, user *models.User) error {
	return r.db.Where("email = ?", email).Preload("Department").First(user).Error
}

func (r *AuthRepository) FindDepartmentByID(id uint) error {
	return r.db.First(&models.Department{}, id).Error
}
