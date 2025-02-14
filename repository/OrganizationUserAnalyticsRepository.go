package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"bbscout/config/database"
	models "bbscout/models/views"
)

type OrganizationUserAnalyticsRepository interface {
	GetOrganizationUserAnalyticsByOrganizationId(organizationId uuid.UUID) ([]*models.OrganizationUserAnalytics, error)
}
type organizationUserAnalyticsRepositoryImpl struct {
	db *gorm.DB
}

func NewOrganizationUserAnalyticsRepository() OrganizationUserAnalyticsRepository {
	if db == nil {
		database.NewDatabaseConnection()
		db = database.GetDB()
	}
	return &organizationUserAnalyticsRepositoryImpl{db: db}
}
func (r *organizationUserAnalyticsRepositoryImpl) GetOrganizationUserAnalyticsByOrganizationId(organizationId uuid.UUID) ([]*models.OrganizationUserAnalytics, error) {
	var organizationAnalytics []*models.OrganizationUserAnalytics
	err := r.db.Where("organization_id = ?", organizationId).Find(&organizationAnalytics).Error
	if err != nil {
		return nil, err
	}
	return organizationAnalytics, nil
}
