package repository

import (
	"bbscout/config/database"
	models "bbscout/models/views"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserMonitoringStatRepository interface {
	GetUserMonitoringStats(organizationId uuid.UUID, userId uuid.UUID) (*models.UserMonitoringStat, error)
}
type userMonitoringStatRepositoryImpl struct {
	db *gorm.DB
}

func NewUserMonitoringStatRepository() UserMonitoringStatRepository {
	if db == nil {
		database.NewDatabaseConnection()
		db = database.GetDB()
	}
	return &userMonitoringStatRepositoryImpl{
		db: db,
	}
}

func (r *userMonitoringStatRepositoryImpl) GetUserMonitoringStats(organizationId uuid.UUID, userId uuid.UUID) (*models.UserMonitoringStat, error) {
	var userMonthlyMonitorStats models.UserMonitoringStat
	err := r.db.Where("organization_id = ? AND user_id = ?", organizationId, userId).First(&userMonthlyMonitorStats).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &userMonthlyMonitorStats, nil
}
