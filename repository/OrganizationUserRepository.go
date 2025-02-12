package repository

import (
	"bbscout/config/database"
	"bbscout/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrganizationUserRepository interface {
	CreateOrganizationUser(organizationUser *models.OrganizationUserModel) (*models.OrganizationUserModel, error)
	GetOrganizationUserById(id uuid.UUID) (*models.OrganizationUserModel, error)
	GetOrganizationUserByUserId(userId uuid.UUID) (*models.OrganizationUserModel, error)
	GetOrganizationUserByOrganizationId(organizationId uuid.UUID) ([]*models.OrganizationUserModel, error)
	UpdateOrganizationUser(organizationUser *models.OrganizationUserModel) (*models.OrganizationUserModel, error)
	DeleteOrganizationUser(id uuid.UUID) error
	ExistUserOrganizationByUserId(userId uuid.UUID, organizationId uuid.UUID) (bool, error)
}
type organizationUserRepositoryImpl struct {
	db *gorm.DB
}

func NewOrganizationUserRepository() OrganizationUserRepository {
	if db == nil {
		database.NewDatabaseConnection()
		db = database.GetDB()
	}
	return &organizationUserRepositoryImpl{db: db}
}
func (r *organizationUserRepositoryImpl) CreateOrganizationUser(organizationUser *models.OrganizationUserModel) (*models.OrganizationUserModel, error) {
	err := r.db.Create(organizationUser).Error
	if err != nil {
		return nil, err
	}
	return organizationUser, nil
}
func (r *organizationUserRepositoryImpl) GetOrganizationUserById(id uuid.UUID) (*models.OrganizationUserModel, error) {
	var organizationUser models.OrganizationUserModel
	err := r.db.First(&organizationUser, id).Error
	if err != nil {
		return nil, err
	}
	return &organizationUser, nil
}
func (r *organizationUserRepositoryImpl) GetOrganizationUserByUserId(userId uuid.UUID) (*models.OrganizationUserModel, error) {
	var organizationUser models.OrganizationUserModel
	err := r.db.Where("user_id = ?", userId).First(&organizationUser).Error
	if err != nil {
		return nil, err
	}
	return &organizationUser, nil
}
func (r *organizationUserRepositoryImpl) GetOrganizationUserByOrganizationId(organizationId uuid.UUID) ([]*models.OrganizationUserModel, error) {
	var organizationUsers []*models.OrganizationUserModel
	err := r.db.Where("organization_id = ?", organizationId).Find(&organizationUsers).Error
	if err != nil {
		return nil, err
	}
	return organizationUsers, nil
}
func (r *organizationUserRepositoryImpl) UpdateOrganizationUser(organizationUser *models.OrganizationUserModel) (*models.OrganizationUserModel, error) {
	err := r.db.Save(organizationUser).Error
	if err != nil {
		return nil, err
	}
	return organizationUser, nil
}
func (r *organizationUserRepositoryImpl) DeleteOrganizationUser(id uuid.UUID) error {
	err := r.db.Delete(&models.OrganizationUserModel{}, id).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *organizationUserRepositoryImpl) ExistUserOrganizationByUserId(userId uuid.UUID, organizationId uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&models.OrganizationUserModel{}).Where("user_id = ? AND organization_id = ?", userId, organizationId).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
