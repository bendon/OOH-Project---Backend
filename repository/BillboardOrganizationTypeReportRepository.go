package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"bbscout/config/database"
	models "bbscout/models/views"
)

type BillboardOrganizationTypeReportRepository interface {
	GetBillboardOrganizationTypeReport(organizationId uuid.UUID, page int, size int, type_ string) ([]models.BillboardOrganizationTypeReport, int64, error)
}

type billboardOrganizationTypeReportRepositoryImpl struct {
	db *gorm.DB
}

func NewBillboardOrganizationTypeReportRepository() BillboardOrganizationTypeReportRepository {
	if db == nil {
		database.NewDatabaseConnection()
		db = database.GetDB()
	}
	return &billboardOrganizationTypeReportRepositoryImpl{
		db: db,
	}
}

func (r *billboardOrganizationTypeReportRepositoryImpl) GetBillboardOrganizationTypeReport(organizationId uuid.UUID, page int, size int, type_ string) ([]models.BillboardOrganizationTypeReport, int64, error) {
	var reports []models.BillboardOrganizationTypeReport
	var total int64
	err := r.db.Order("type_count DESC").Where("organization_id = ? AND type LIKE ?", organizationId, "%"+type_+"%").Offset((page - 1) * size).Limit(size).Find(&reports).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	return reports, total, nil
}
