package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"bbscout/config/database"
	models "bbscout/models/views"
)

type BillboardOrganizationLocationReportRepository interface {
	GetBillboardOrganizationLocationReport(organizationId uuid.UUID, page int, size int, location string) ([]models.BillboardOrganizationLocationReport, int64, error)
}
type billboardOrganizationLocationReportRepositoryImpl struct {
	db *gorm.DB
}

func NewBillboardOrganizationLocationReportRepository() BillboardOrganizationLocationReportRepository {
	if db == nil {
		database.NewDatabaseConnection()
		db = database.GetDB()
	}
	return &billboardOrganizationLocationReportRepositoryImpl{
		db: db,
	}
}
func (r *billboardOrganizationLocationReportRepositoryImpl) GetBillboardOrganizationLocationReport(organizationId uuid.UUID, page int, size int, location string) ([]models.BillboardOrganizationLocationReport, int64, error) {
	var reports []models.BillboardOrganizationLocationReport
	var total int64
	err := r.db.Order("count_per_location DESC").Where("organization_id = ? AND location LIKE ?", organizationId, "%"+location+"%").Offset((page - 1) * size).Limit(size).Find(&reports).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	return reports, total, nil
}
