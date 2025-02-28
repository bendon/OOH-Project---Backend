package repository

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"bbscout/config/database"
	models "bbscout/models/views"
)

type UserBillboardUploadReportRepository interface {
	GetUserBillboardUploadReportByUser(organizationId uuid.UUID, userID uuid.UUID) (*models.UserBillboardUploadReport, error)
	GetUserBillboardUploadReportByOrganizationPageable(organizationId uuid.UUID, page int, size int, user string) ([]models.UserBillboardUploadReport, int64, error)
}
type userBillboardUploadReportRepositoryImpl struct {
	db *gorm.DB
}

func NewUserBillboardUploadReportRepository() UserBillboardUploadReportRepository {
	if db == nil {
		database.NewDatabaseConnection()
		db = database.GetDB()
	}
	return &userBillboardUploadReportRepositoryImpl{
		db: db,
	}
}

func (r *userBillboardUploadReportRepositoryImpl) GetUserBillboardUploadReportByUser(organizationId uuid.UUID, userID uuid.UUID) (*models.UserBillboardUploadReport, error) {
	var report models.UserBillboardUploadReport
	err := r.db.Preload("UserInfo").Where("organization_id = ? AND user_id = ?", organizationId, userID).First(&report).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil

		}
		return nil, err
	}
	return &report, nil
}

func (r *userBillboardUploadReportRepositoryImpl) GetUserBillboardUploadReportByOrganizationPageable(organizationId uuid.UUID, page int, size int, user string) ([]models.UserBillboardUploadReport, int64, error) {
	var reports []models.UserBillboardUploadReport
	var total int64
	err := r.db.Preload("UserInfo").Where("organization_id = ? AND user_name LIKE ?", organizationId, "%"+user+"%").Offset((page - 1) * size).Limit(size).Find(&reports).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	return reports, total, nil

}
