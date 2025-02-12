package repository

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"bbscout/config/database"
	"bbscout/models"
)

type UserRoleRepository interface {
	CreateUserRole(userRole *models.UserRoleModel) (*models.UserRoleModel, error)
	GetUserRoleById(id uuid.UUID) (*models.UserRoleModel, error)
	UpdateUserRole(userRole *models.UserRoleModel) (*models.UserRoleModel, error)
	DeleteUserRole(id uuid.UUID) error
	GetUserRoleByUserIdAndOrganizationId(userId uuid.UUID, organizationId uuid.UUID) (*models.UserRoleModel, error)
}
type userRoleRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRoleRepository() UserRoleRepository {
	if db == nil {
		database.NewDatabaseConnection()
		db = database.GetDB()
	}
	return &userRoleRepositoryImpl{db: db}
}
func (r *userRoleRepositoryImpl) CreateUserRole(userRole *models.UserRoleModel) (*models.UserRoleModel, error) {
	err := r.db.Create(userRole).Error
	if err != nil {
		return nil, err
	}
	return userRole, nil
}
func (r *userRoleRepositoryImpl) GetUserRoleById(id uuid.UUID) (*models.UserRoleModel, error) {
	var userRole models.UserRoleModel
	err := r.db.First(&userRole, id).Error
	if err != nil {
		return nil, err
	}
	return &userRole, nil
}
func (r *userRoleRepositoryImpl) UpdateUserRole(userRole *models.UserRoleModel) (*models.UserRoleModel, error) {
	err := r.db.Save(userRole).Error
	if err != nil {
		return nil, err
	}
	return userRole, nil
}
func (r *userRoleRepositoryImpl) DeleteUserRole(id uuid.UUID) error {
	err := r.db.Delete(&models.UserRoleModel{}, id).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *userRoleRepositoryImpl) GetUserRoleByUserIdAndOrganizationId(userId uuid.UUID, organizationId uuid.UUID) (*models.UserRoleModel, error) {
	var userRole models.UserRoleModel
	err := r.db.Where("user_id = ? AND organization_id = ?", userId, organizationId).First(&userRole).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil

		}
		return nil, err
	}
	return &userRole, nil
}
