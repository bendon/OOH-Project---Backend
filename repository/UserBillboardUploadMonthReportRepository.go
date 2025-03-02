package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"bbscout/config/database"
	models "bbscout/models/views"
)

type UserBillboardUploadMonthReportRepository interface {
	GetUserBillboardUploadMonthReportByUserMonthyly(organizationId uuid.UUID, userID uuid.UUID, year int, month int) (*models.UserBillboardUploadMonthReport, error)
	GetUserBillboardUploadMonthReportByUserYearly(organizationId uuid.UUID, userID uuid.UUID, year int) ([]models.UserBillboardUploadMonthReport, error)
}
type userBillboardUploadMonthReportRepositoryImpl struct {
	db *gorm.DB
}

func NewUserBillboardUploadMonthReportRepository() UserBillboardUploadMonthReportRepository {
	if db == nil {
		database.NewDatabaseConnection()
		db = database.GetDB()
	}
	return &userBillboardUploadMonthReportRepositoryImpl{
		db: db,
	}
}

func (r *userBillboardUploadMonthReportRepositoryImpl) GetUserBillboardUploadMonthReportByUserYearly(organizationId uuid.UUID, userID uuid.UUID, year int) ([]models.UserBillboardUploadMonthReport, error) {
	var reports []models.UserBillboardUploadMonthReport
	err := r.db.Where("organization_id = ? AND user_id = ? AND upload_year = ?", organizationId, userID, year).Find(&reports).Error
	if err != nil {
		return nil, err
	}
	return reports, nil
}

func (r *userBillboardUploadMonthReportRepositoryImpl) GetUserBillboardUploadMonthReportByUserMonthyly(organizationId uuid.UUID, userID uuid.UUID, year int, month int) (*models.UserBillboardUploadMonthReport, error) {
	var report models.UserBillboardUploadMonthReport
	err := r.db.Where("organization_id = ? AND user_id = ? AND upload_year = ? AND upload_month = ?", organizationId, userID, year, month).First(&report).Error
	if err != nil {
		return nil, err
	}
	return &report, nil
}
