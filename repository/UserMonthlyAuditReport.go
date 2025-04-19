package repository

import (
	"bbscout/config/database"
	models "bbscout/models/views"
	"bbscout/types"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserMonthlyAuditReportRepository interface {
	GetUserMonthlyAuditReport(organizationId uuid.UUID, userId uuid.UUID, year int) ([]types.MonthlyAuditReport, error)
}
type userMonthlyAuditReportRepositoryImpl struct {
	db *gorm.DB
}

func NewUserMonthlyAuditReportRepository() UserMonthlyAuditReportRepository {
	if db == nil {
		database.NewDatabaseConnection()
		db = database.GetDB()
	}
	return &userMonthlyAuditReportRepositoryImpl{
		db: db,
	}
}
func (r *userMonthlyAuditReportRepositoryImpl) GetUserMonthlyAuditReport(organizationId uuid.UUID, userId uuid.UUID, year int) ([]types.MonthlyAuditReport, error) {
	var userMonthlyAuditStats []models.UserMonthlyAuditStat
	err := r.db.Where("user_id = ? AND organization_id = ? AND year_number = ?", userId, organizationId, year).Find(&userMonthlyAuditStats).Error
	if err != nil {
		return nil, err
	}

	reports := make([]types.MonthlyAuditReport, len(userMonthlyAuditStats))
	for i, stat := range userMonthlyAuditStats {
		reports[i] = types.MonthlyAuditReport{
			Month:      stat.MonthNumber,
			MonthName:  stat.MonthName,
			TotalAudit: stat.TotalBillboards,
		}
	}

	return reports, nil
}
