package repository

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"bbscout/config/database"
	"bbscout/models"
)

type RoleRepository interface {
	CreateRole(role *models.RoleModel) (*models.RoleModel, error)
	GetRoleById(id uuid.UUID) (*models.RoleModel, error)
	UpdateRole(role *models.RoleModel) (*models.RoleModel, error)
	GetRoleByName(name string) (*models.RoleModel, error)
	DeleteRole(id uuid.UUID) error
	GetRolesByOrganizationId(organizationId uuid.UUID) ([]models.RoleModel, error)
	ExistsRoleByIdAndOrganizationId(id uuid.UUID, organizationId uuid.UUID) (bool, error)
	ExistsRoleByNameAndOrganizationId(name string, organizationId uuid.UUID) (bool, error)
	GetRoleByIdAndOrganizationId(id uuid.UUID, organizationId uuid.UUID) (*models.RoleModel, error)
	GetRoleByNameAndOrganizationId(name string, organizationId uuid.UUID) (*models.RoleModel, error)
}
type roleRepositoryImpl struct {
	db *gorm.DB
}

func NewRoleRepository() RoleRepository {
	if db == nil {
		database.NewDatabaseConnection()
		db = database.GetDB()
	}
	return &roleRepositoryImpl{db: db}
}
func (r *roleRepositoryImpl) CreateRole(role *models.RoleModel) (*models.RoleModel, error) {
	err := r.db.FirstOrCreate(&models.RoleModel{}, role).Error
	if err != nil {
		return nil, err
	}
	return role, nil
}
func (r *roleRepositoryImpl) GetRoleById(id uuid.UUID) (*models.RoleModel, error) {
	var role models.RoleModel
	err := r.db.First(&role, id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepositoryImpl) UpdateRole(role *models.RoleModel) (*models.RoleModel, error) {
	err := r.db.Save(role).Error
	if err != nil {
		return nil, err
	}
	return role, nil
}
func (r *roleRepositoryImpl) DeleteRole(id uuid.UUID) error {
	err := r.db.Delete(&models.RoleModel{}, id).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *roleRepositoryImpl) GetRoleByName(name string) (*models.RoleModel, error) {
	var role models.RoleModel
	err := r.db.Where("name = ?", name).First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil

		}
		return nil, err
	}
	return &role, nil
}

func (r *roleRepositoryImpl) GetRolesByOrganizationId(organizationId uuid.UUID) ([]models.RoleModel, error) {
	var roles []models.RoleModel
	err := r.db.Where("organization_id = ?", organizationId).Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}
func (r *roleRepositoryImpl) ExistsRoleByIdAndOrganizationId(id uuid.UUID, organizationId uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&models.RoleModel{}).Where("id = ? AND organization_id = ?", id, organizationId).Count(&count).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil

		}
		return false, err
	}
	return count > 0, nil
}
func (r *roleRepositoryImpl) ExistsRoleByNameAndOrganizationId(name string, organizationId uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&models.RoleModel{}).Where("name = ? AND organization_id = ?", name, organizationId).Count(&count).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil

		}
		return false, err
	}
	return count > 0, nil
}
func (r *roleRepositoryImpl) GetRoleByIdAndOrganizationId(id uuid.UUID, organizationId uuid.UUID) (*models.RoleModel, error) {
	var role models.RoleModel
	err := r.db.Where("id = ? AND organization_id = ?", id, organizationId).First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil

		}
		return nil, err
	}
	return &role, nil
}

func (r *roleRepositoryImpl) GetRoleByNameAndOrganizationId(name string, organizationId uuid.UUID) (*models.RoleModel, error) {
	var role models.RoleModel
	err := r.db.Where("name = ? AND organization_id = ?", name, organizationId).First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}
