package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"bbscout/config/database"
	models "bbscout/models/views"
)

type UserBillboardUploadsReportWeekRepository interface {
	GetUserBillboardUploadsReportWeekByUser(organizationId uuid.UUID, userId uuid.UUID, year int, month int, week int) ([]models.UserBillboardUploadsReportWeek, error)
}
type userBillboardUploadsReportWeekRepositoryImpl struct {
	db *gorm.DB
}

func NewUserBillboardUploadsReportWeekRepository() UserBillboardUploadsReportWeekRepository {
	if db == nil {
		database.NewDatabaseConnection()
		db = database.GetDB()
	}
	return &userBillboardUploadsReportWeekRepositoryImpl{
		db: db,
	}
}

func (r *userBillboardUploadsReportWeekRepositoryImpl) GetUserBillboardUploadsReportWeekByUser(organizationId uuid.UUID, userId uuid.UUID, year int, month int, week int) ([]models.UserBillboardUploadsReportWeek, error) {
	var reports []models.UserBillboardUploadsReportWeek
	err := r.db.Where("organization_id = ? AND user_id = ? AND upload_year = ? AND upload_month = ? AND week_number = ?", organizationId, userId, year, month, week).Find(&reports).Error
	if err != nil {
		return nil, err
	}
	return reports, nil
}
