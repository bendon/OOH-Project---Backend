package repository

import (
	"bbscout/config/database"
	models "bbscout/models/views"
	"bbscout/types"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserMonthlyMonitorStatRepository interface {
	GetUserMonthlyMonitorStats(organizationId uuid.UUID, userId uuid.UUID, year int) ([]types.MonthlyMonitorReport, error)
}

type userMonthlyMonitorStatRepositoryImpl struct {
	db *gorm.DB
}

func NewUserMonthlyMonitorStatRepository() UserMonthlyMonitorStatRepository {
	if db == nil {
		database.NewDatabaseConnection()
		db = database.GetDB()
	}
	return &userMonthlyMonitorStatRepositoryImpl{
		db: db,
	}
}

func (r *userMonthlyMonitorStatRepositoryImpl) GetUserMonthlyMonitorStats(organizationId uuid.UUID, userId uuid.UUID, year int) ([]types.MonthlyMonitorReport, error) {
	var userMonthlyMonitorStats []models.UserMonthlyMonitoringStat
	err := r.db.Where("organization_id = ? AND user_id = ? AND year_number = ?", organizationId, userId, year).Find(&userMonthlyMonitorStats).Error
	if err != nil {
		return nil, err
	}

	reports := make([]types.MonthlyMonitorReport, len(userMonthlyMonitorStats))
	for i, stat := range userMonthlyMonitorStats {
		reports[i] = types.MonthlyMonitorReport{
			Month:        stat.MonthNumber,
			MonthName:    stat.MonthName,
			TotalMonitor: stat.TotalMonitors,
		}
	}

	return reports, nil
}
