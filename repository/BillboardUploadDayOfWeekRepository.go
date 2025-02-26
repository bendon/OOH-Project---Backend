package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"bbscout/config/database"
	models "bbscout/models/views"
)

type BillboardUploadDayOfWeekRepository interface {
	GetBillboardUploadDayOfWeek(organizationId uuid.UUID, year int, month int, weekNumber int) ([]models.BillboardUploadDayOfWeek, error)
}

type billboardUploadDayOfWeekRepositoryImpl struct {
	db *gorm.DB
}

func NewBillboardUploadDayOfWeekRepository() BillboardUploadDayOfWeekRepository {
	if db == nil {
		database.NewDatabaseConnection()
		db = database.GetDB()
	}
	return &billboardUploadDayOfWeekRepositoryImpl{
		db: db,
	}
}
func (r *billboardUploadDayOfWeekRepositoryImpl) GetBillboardUploadDayOfWeek(organizationId uuid.UUID, year int, month int, weekNumber int) ([]models.BillboardUploadDayOfWeek, error) {
	var billboards []models.BillboardUploadDayOfWeek
	err := r.db.Where("organization_id = ? AND upload_year = ? AND upload_month = ? AND upload_week_number = ?", organizationId, year, month, weekNumber).Find(&billboards).Error
	if err != nil {
		return nil, err
	}
	return billboards, nil
}
