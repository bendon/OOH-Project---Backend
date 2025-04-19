package repository

import (
	"bbscout/config/database"
	models "bbscout/models/views"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserAuditReportRepository interface {
	GetUserAuditReport(organizationId uuid.UUID, userId uuid.UUID) (*models.UserAuditReport, error)
}

type userAuditReportRepositoryImpl struct {
	db *gorm.DB
}

func NewUserAuditReportRepository() UserAuditReportRepository {
	if db == nil {
		database.NewDatabaseConnection()
		db = database.GetDB()
	}
	return &userAuditReportRepositoryImpl{
		db: db,
	}
}
func (r *userAuditReportRepositoryImpl) GetUserAuditReport(organizationId uuid.UUID, userId uuid.UUID) (*models.UserAuditReport, error) {
	var userAuditReport models.UserAuditReport
	err := r.db.Where("user_id = ? AND organization_id = ?", userId, organizationId).First(&userAuditReport).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &userAuditReport, nil
}
