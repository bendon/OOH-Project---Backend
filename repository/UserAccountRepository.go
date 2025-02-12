package repository

import (
	"bbscout/config/database"
	"bbscout/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserAccountRepository interface {
	CreateUserAccount(userAccount *models.UserAccountModel) (*models.UserAccountModel, error)
	GetUserAccountById(id uuid.UUID) (*models.UserAccountModel, error)
	UpdateUserAccount(userAccount *models.UserAccountModel) (*models.UserAccountModel, error)
	DeleteUserAccount(id uuid.UUID) error
	GetUserAccountsByUserId(userId uuid.UUID) ([]models.UserAccountModel, error)
	GetUserAccountByUserIdAndId(userId uuid.UUID, id uuid.UUID) (*models.UserAccountModel, error)
	GetUserAccountsByIdAndOrganizationId(id uuid.UUID, organizationId uuid.UUID) (*models.UserAccountModel, error)
	GetUserAccountByUserIdAndOrganizationId(userId uuid.UUID, organizationId uuid.UUID) (*models.UserAccountModel, error)
}
type userAccountRepositoryImpl struct {
	db *gorm.DB
}

func NewUserAccountRepository() UserAccountRepository {
	if db == nil {
		database.NewDatabaseConnection()
		db = database.GetDB()
	}
	return &userAccountRepositoryImpl{db: db}
}
func (r *userAccountRepositoryImpl) CreateUserAccount(userAccount *models.UserAccountModel) (*models.UserAccountModel, error) {
	err := r.db.Create(userAccount).Error
	if err != nil {
		return nil, err
	}
	return userAccount, nil
}
func (r *userAccountRepositoryImpl) GetUserAccountById(id uuid.UUID) (*models.UserAccountModel, error) {
	var userAccount models.UserAccountModel
	err := r.db.Preload("Organization").First(&userAccount, id).Error
	if err != nil {
		return nil, err
	}
	return &userAccount, nil

}
func (r *userAccountRepositoryImpl) UpdateUserAccount(userAccount *models.UserAccountModel) (*models.UserAccountModel, error) {
	err := r.db.Save(userAccount).Error
	if err != nil {
		return nil, err
	}
	return userAccount, nil
}

func (r *userAccountRepositoryImpl) DeleteUserAccount(id uuid.UUID) error {
	var userAccount models.UserAccountModel
	err := r.db.Delete(&userAccount, id).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *userAccountRepositoryImpl) GetUserAccountsByUserId(userId uuid.UUID) ([]models.UserAccountModel, error) {
	var userAccounts []models.UserAccountModel
	err := r.db.Preload("Organization").Where("user_id = ?", userId).Where("active=?", true).Where("is_locked=?", false).Find(&userAccounts).Error
	if err != nil {
		return nil, err
	}
	return userAccounts, nil
}
func (r *userAccountRepositoryImpl) GetUserAccountByUserIdAndId(userId uuid.UUID, id uuid.UUID) (*models.UserAccountModel, error) {
	var userAccount models.UserAccountModel
	err := r.db.Preload("Organization").Where("user_id = ?", userId).Where("id = ?", id).Where("active=?", true).Where("is_locked=?", false).First(&userAccount).Error
	if err != nil {
		return nil, err
	}
	return &userAccount, nil
}
func (r *userAccountRepositoryImpl) GetUserAccountsByIdAndOrganizationId(id uuid.UUID, organizationId uuid.UUID) (*models.UserAccountModel, error) {
	var userAccount models.UserAccountModel
	err := r.db.Preload("Organization").Where("id = ?", id).Where("organization_id = ?", organizationId).First(&userAccount).Error
	if err != nil {
		return nil, err
	}
	return &userAccount, nil
}
func (r *userAccountRepositoryImpl) GetUserAccountByUserIdAndOrganizationId(userId uuid.UUID, organizationId uuid.UUID) (*models.UserAccountModel, error) {
	var userAccount models.UserAccountModel
	err := r.db.Preload("Organization").Where("user_id = ?", userId).Where("organization_id = ?", organizationId).First(&userAccount).Error
	if err != nil {
		return nil, err
	}
	return &userAccount, nil
}
