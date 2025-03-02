package repository

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"bbscout/config/database"
	"bbscout/models"
)

type OrganizationRepository interface {
	CreateOrganization(organization *models.OrganizationModel) (*models.OrganizationModel, error)
	GetOrganizationById(id uuid.UUID) (*models.OrganizationModel, error)
	UpdateOrganization(organization *models.OrganizationModel) (*models.OrganizationModel, error)
	GetOperationOrganization() (*models.OrganizationModel, error)
	DeleteOrganization(id uuid.UUID) error
	GetOrganizationDetailsById(id uuid.UUID) (*models.OrganizationModel, error)
	FindOrganizationLatest() (*models.OrganizationModel, error)
}
type organizationRepositoryImpl struct {
	db *gorm.DB
}

func NewOrganizationRepository() OrganizationRepository {
	if db == nil {
		database.NewDatabaseConnection()
		db = database.GetDB()
	}
	return &organizationRepositoryImpl{db: db}
}
func (r *organizationRepositoryImpl) CreateOrganization(organization *models.OrganizationModel) (*models.OrganizationModel, error) {
	err := r.db.Create(organization).Error
	if err != nil {
		return nil, err
	}
	return organization, nil
}
func (r *organizationRepositoryImpl) GetOrganizationById(organizationId uuid.UUID) (*models.OrganizationModel, error) {
	var organization models.OrganizationModel
	err := r.db.First(&organization, organizationId)
	if err.Error != nil {
		if errors.Is(err.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err.Error
	}
	return &organization, nil
}
func (r *organizationRepositoryImpl) UpdateOrganization(organization *models.OrganizationModel) (*models.OrganizationModel, error) {
	err := r.db.Save(organization).Error
	if err != nil {
		return nil, err
	}
	return organization, nil
}
func (r *organizationRepositoryImpl) DeleteOrganization(organizationId uuid.UUID) error {
	err := r.db.Delete(&models.OrganizationModel{}, organizationId)
	if err.Error != nil {
		if errors.Is(err.Error, gorm.ErrRecordNotFound) {
			return nil
		}
		return err.Error
	}
	return nil
}
func (r *organizationRepositoryImpl) GetOperationOrganization() (*models.OrganizationModel, error) {
	var operationOrganization models.OrganizationModel
	err := r.db.Where("is_operation = ?", true).First(&operationOrganization).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &operationOrganization, nil
}
func (r *organizationRepositoryImpl) GetOrganizationDetailsById(id uuid.UUID) (*models.OrganizationModel, error) {
	var organization models.OrganizationModel
	err := r.db.Preload("Admin").First(&organization, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &organization, nil
}

func (r *organizationRepositoryImpl) FindOrganizationLatest() (*models.OrganizationModel, error) {
	var organization models.OrganizationModel
	err := r.db.Last(&organization).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &organization, nil
}
