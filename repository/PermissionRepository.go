package repository

import (
	"bbscout/config/database"
	"bbscout/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PermissionRepository interface {
	CreatePermission(permission *models.PermissionModel) (*models.PermissionModel, error)
	GetPermissionById(id uuid.UUID) (*models.PermissionModel, error)
	UpdatePermission(permission *models.PermissionModel) (*models.PermissionModel, error)
	GetPermissionsByAccount(account string) ([]models.PermissionModel, error)
	DeletePermission(id uuid.UUID) error
	ExistsPermissionById(id uuid.UUID) (bool, error)
	GetPermissionsByNotType(account string) ([]models.PermissionModel, error)
}
type permissionRepositoryImpl struct {
	db *gorm.DB
}

func NewPermissionRepository() PermissionRepository {
	if db == nil {
		database.NewDatabaseConnection()
		db = database.GetDB()
	}
	return &permissionRepositoryImpl{db: db}
}
func (r *permissionRepositoryImpl) CreatePermission(permission *models.PermissionModel) (*models.PermissionModel, error) {
	// Check for existing permission by name, type, and account
	var existingPermission models.PermissionModel
	err := r.db.Where("name = ? AND type = ? AND account = ?", permission.Name, permission.Type, permission.Account).First(&existingPermission).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// If no existing permission, create a new one
	if err == gorm.ErrRecordNotFound {
		err = r.db.Create(permission).Error
		if err != nil {
			return nil, err
		}
	}
	return permission, nil

}
func (r *permissionRepositoryImpl) GetPermissionById(id uuid.UUID) (*models.PermissionModel, error) {
	var permission models.PermissionModel
	err := r.db.First(&permission, id).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}
func (r *permissionRepositoryImpl) UpdatePermission(permission *models.PermissionModel) (*models.PermissionModel, error) {
	err := r.db.Save(permission).Error
	if err != nil {
		return nil, err
	}
	return permission, nil
}
func (r *permissionRepositoryImpl) DeletePermission(id uuid.UUID) error {
	err := r.db.Delete(&models.PermissionModel{}, id).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *permissionRepositoryImpl) GetPermissionsByAccount(account string) ([]models.PermissionModel, error) {
	var permissions []models.PermissionModel
	err := r.db.Where("account = ?", account).Find(&permissions).Error
	if err != nil {
		return nil, err
	}
	return permissions, nil
}
func (r *permissionRepositoryImpl) ExistsPermissionById(id uuid.UUID) (bool, error) {
	var permission models.PermissionModel
	err := r.db.Where("id = ?", id).First(&permission).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
func (r *permissionRepositoryImpl) GetPermissionsByNotType(account string) ([]models.PermissionModel, error) {
	var permissions []models.PermissionModel
	err := r.db.Where("type != ?", account).Find(&permissions).Error
	if err != nil {
		return nil, err
	}
	return permissions, nil
}
