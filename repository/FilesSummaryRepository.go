package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"bbscout/config/database"
	models "bbscout/models/views"
)

type FilesSummaryRepository interface {
	GetFilesSummaryPageable(organizationId uuid.UUID, startDate int64, endDate int64, page int, size int) ([]models.FilesSummary, int64, error)
}
type filesSummaryRepositoryImpl struct {
	db *gorm.DB
}

func NewFilesSummaryRepository() FilesSummaryRepository {
	if db == nil {
		database.NewDatabaseConnection()
		db = database.GetDB()
	}
	return &filesSummaryRepositoryImpl{
		db: db,
	}
}

func (r *filesSummaryRepositoryImpl) GetFilesSummaryPageable(organizationId uuid.UUID, startDate int64, endDate int64, page int, size int) ([]models.FilesSummary, int64, error) {
	var files []models.FilesSummary
	var totalCount int64
	if err := r.db.Where("organization_id = ? AND created_at BETWEEN ? AND ?", organizationId, startDate, endDate).Model(&models.FilesSummary{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	if err := r.db.Preload("Billboard").Where("organization_id = ? AND created_at BETWEEN ? AND ?", organizationId, startDate, endDate).
		Limit(size).
		Offset(offset).
		Find(&files).Error; err != nil {
		return nil, 0, err
	}
	return files, totalCount, nil
}
