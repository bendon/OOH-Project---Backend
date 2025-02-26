package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"bbscout/config/database"
	models "bbscout/models/views"
)

type BillboardOrganizationReportRepository interface {
	GetBillboardOrganizationReport(organizationId uuid.UUID) (*models.BillboardOrganizationReport, error)
}

type billboardOrganizationReportRepositoryImpl struct {
	db *gorm.DB
}

func NewBillboardOrganizationReportRepository() BillboardOrganizationReportRepository {
	if db == nil {
		database.NewDatabaseConnection()
		db = database.GetDB()
	}
	return &billboardOrganizationReportRepositoryImpl{
		db: db,
	}
}
func (r *billboardOrganizationReportRepositoryImpl) GetBillboardOrganizationReport(organizationId uuid.UUID) (*models.BillboardOrganizationReport, error) {
	var report models.BillboardOrganizationReport
	err := r.db.Where("organization_id = ?", organizationId).First(&report).Error
	if err != nil {
		return nil, err
	}
	return &report, nil
}
