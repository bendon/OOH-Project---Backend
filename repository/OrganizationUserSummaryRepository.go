package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"bbscout/config/database"
	models "bbscout/models/views"
)

type OrganizationUserSummaryRepository interface {
	GetOrganizationUserSummary(organizationID uuid.UUID) ([]models.OrganizationUserSummary, error)
}

type organizationUserSummaryRepositoryImpl struct {
	db *gorm.DB
}

func NewOrganizationUserSummaryRepository() OrganizationUserSummaryRepository {
	if db == nil {
		database.NewDatabaseConnection()
		db = database.GetDB()
	}
	return &organizationUserSummaryRepositoryImpl{
		db: db,
	}
}

func (r *organizationUserSummaryRepositoryImpl) GetOrganizationUserSummary(organizationID uuid.UUID) ([]models.OrganizationUserSummary, error) {
	var organizationUserSummary []models.OrganizationUserSummary
	err := r.db.Where("organization_id = ?", organizationID).Find(&organizationUserSummary).Error
	if err != nil {
		return nil, err
	}
	return organizationUserSummary, nil
}
