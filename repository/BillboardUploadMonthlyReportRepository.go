package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"bbscout/config/database"
	models "bbscout/models/views"
)

type BillboardUploadMonthlyReportRepository interface {
	GetBillboardUploadMonthlyReport(organizationId uuid.UUID, year int, month int) ([]models.BillboardUploadMonthlyReport, error)
	GetBillboardUploadMonthlyReportByYear(organizationId uuid.UUID, year int) ([]models.BillboardUploadMonthlyReport, error)
}
type billboardUploadMonthlyReportRepositoryImpl struct {
	db *gorm.DB
}

func NewBillboardUploadMonthlyReportRepository() BillboardUploadMonthlyReportRepository {
	if db == nil {
		database.NewDatabaseConnection()
		db = database.GetDB()
	}
	return &billboardUploadMonthlyReportRepositoryImpl{
		db: db,
	}
}

func (r *billboardUploadMonthlyReportRepositoryImpl) GetBillboardUploadMonthlyReport(organizationId uuid.UUID, year int, month int) ([]models.BillboardUploadMonthlyReport, error) {
	var billboards []models.BillboardUploadMonthlyReport
	err := r.db.Where("organization_id = ? AND upload_year = ? AND upload_month = ?", organizationId, year, month).Find(&billboards).Error
	if err != nil {
		return nil, err
	}
	return billboards, nil
}
func (r *billboardUploadMonthlyReportRepositoryImpl) GetBillboardUploadMonthlyReportByYear(organizationId uuid.UUID, year int) ([]models.BillboardUploadMonthlyReport, error) {
	var billboards []models.BillboardUploadMonthlyReport
	err := r.db.Where("organization_id = ? AND upload_year = ?", organizationId, year).Find(&billboards).Error
	if err != nil {
		return nil, err
	}
	return billboards, nil
}
