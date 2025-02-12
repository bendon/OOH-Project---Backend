package repository

import (
	"bbscout/config/database"
	"bbscout/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserAccountPermissionRepository interface {
	CreateUserAccountPermission(userAccountPermission *models.UserAccountPermissionModel) (*models.UserAccountPermissionModel, error)
	GetUserAccountPermissionById(id uuid.UUID) (*models.UserAccountPermissionModel, error)
	UpdateUserAccountPermission(userAccountPermission *models.UserAccountPermissionModel) (*models.UserAccountPermissionModel, error)
	DeleteUserAccountPermission(id uuid.UUID) error
	GetPermissionsByUserIdAndAccountId(userId uuid.UUID, accountId uuid.UUID) ([]models.PermissionModel, error)
	GetUserAccountPermissionsByOrganizationId(accountId uuid.UUID, organizationId uuid.UUID) ([]models.PermissionModel, error)
	DeleteUserAccountPermissionsByAccountId(accountId uuid.UUID) error
	ExistsUserAccountPermissionByAccountIdAndPermissionId(accountId uuid.UUID, permissionId uuid.UUID, organizationId uuid.UUID) (bool, error)
}
type userAccountPermissionRepositoryImpl struct {
	db *gorm.DB
}

func NewUserAccountPermissionRepository() UserAccountPermissionRepository {
	if db == nil {
		database.NewDatabaseConnection()
		db = database.GetDB()
	}
	return &userAccountPermissionRepositoryImpl{db: db}
}
func (r *userAccountPermissionRepositoryImpl) CreateUserAccountPermission(userAccountPermission *models.UserAccountPermissionModel) (*models.UserAccountPermissionModel, error) {
	err := r.db.Create(userAccountPermission).Error
	if err != nil {
		return nil, err
	}
	return userAccountPermission, nil
}
func (r *userAccountPermissionRepositoryImpl) GetUserAccountPermissionById(id uuid.UUID) (*models.UserAccountPermissionModel, error) {
	var userAccountPermission models.UserAccountPermissionModel
	err := r.db.First(&userAccountPermission, id).Error
	if err != nil {
		return nil, err
	}
	return &userAccountPermission, nil
}
func (r *userAccountPermissionRepositoryImpl) UpdateUserAccountPermission(userAccountPermission *models.UserAccountPermissionModel) (*models.UserAccountPermissionModel, error) {
	err := r.db.Save(userAccountPermission).Error
	if err != nil {
		return nil, err
	}
	return userAccountPermission, nil
}
func (r *userAccountPermissionRepositoryImpl) DeleteUserAccountPermission(id uuid.UUID) error {
	err := r.db.Delete(&models.UserAccountPermissionModel{}, id).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *userAccountPermissionRepositoryImpl) GetPermissionsByUserIdAndAccountId(userId uuid.UUID, accountId uuid.UUID) ([]models.PermissionModel, error) {
	var permissions []models.PermissionModel

	// Preload Branches from BranchUsersModel, but return only the Branches
	if err := db.Model(&models.UserAccountPermissionModel{}).
		Select("permissions.*").
		Joins("JOIN permissions ON permissions.id = user_account_permissions.permission_id").
		Where("user_account_permissions.user_id = ? AND user_account_permissions.account_id = ?", userId, accountId).
		Find(&permissions).Error; err != nil {
		return nil, err
	}

	return permissions, nil

}
func (r *userAccountPermissionRepositoryImpl) GetUserAccountPermissionsByOrganizationId(accountId uuid.UUID, organizationId uuid.UUID) ([]models.PermissionModel, error) {

	var permissions []models.PermissionModel

	// Preload Branches from BranchUsersModel, but return only the Branches
	if err := db.Model(&models.UserAccountPermissionModel{}).
		Select("permissions.*").
		Joins("JOIN permissions ON permissions.id = user_account_permissions.permission_id").
		Where("user_account_permissions.organization_id = ? AND user_account_permissions.account_id = ?", organizationId, accountId).
		Find(&permissions).Error; err != nil {
		return nil, err
	}

	return permissions, nil
}

func (r *userAccountPermissionRepositoryImpl) DeleteUserAccountPermissionsByAccountId(accountId uuid.UUID) error {
	err := r.db.Where("account_id = ?", accountId).Delete(&models.UserAccountPermissionModel{}).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *userAccountPermissionRepositoryImpl) ExistsUserAccountPermissionByAccountIdAndPermissionId(accountId uuid.UUID, permissionId uuid.UUID, organizationId uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&models.UserAccountPermissionModel{}).
		Where("account_id = ? AND permission_id = ? AND organization_id = ?", accountId, permissionId, organizationId).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
