package repository

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"bbscout/config/database"
	"bbscout/models"
)

var db *gorm.DB

type UserRepository interface {
	CreateUser(user *models.UserModel) (*models.UserModel, error)
	GetUserById(id uuid.UUID) (*models.UserModel, error)
	UpdateUser(user *models.UserModel) (*models.UserModel, error)
	DeleteUser(id uuid.UUID) error
	GetUserByEmail(email string) (*models.UserModel, error)
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

func (r *userRepositoryImpl) GetUserById(id uuid.UUID) (*models.UserModel, error) {
	var user models.UserModel
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func (r *userRepositoryImpl) UpdateUser(user *models.UserModel) (*models.UserModel, error) {
	err := r.db.Save(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (r *userRepositoryImpl) DeleteUser(id uuid.UUID) error {
	err := r.db.Delete(&models.UserModel{}, id).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *userRepositoryImpl) GetUserByEmail(email string) (*models.UserModel, error) {
	var user models.UserModel
	err := r.db.Where("email = ? ", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // User not found
		}
		return nil, err
	}
	return &user, nil
}
