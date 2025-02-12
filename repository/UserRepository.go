package repository

import (
	"gorm.io/gorm"

	"bbscout/config/database"
	"bbscout/models"
)

var db *gorm.DB

type UserRepository interface {
	CreateUser(user *models.UserModel) (*models.UserModel, error)
}

type userRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository() UserRepository {
	if db == nil {
		database.NewDatabaseConnection()
		db = database.GetDB()
	}
	return &userRepositoryImpl{
		db: db,
	}
}

func (r *userRepositoryImpl) CreateUser(user *models.UserModel) (*models.UserModel, error) {
	err := r.db.Create(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}
