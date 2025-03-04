package repository

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"bbscout/config/database"
	models "bbscout/models/views"
)

type BillboardSummaryRepository interface {
	GetStaffBillBoardsSummary(organizationId uuid.UUID, createdById uuid.UUID, page int, size int, search string) ([]models.BillboardSummaryView, int64, error)
	GetStaffBillBoardsSummaryById(id uuid.UUID, createdById uuid.UUID) (*models.BillboardSummaryView, error)
	GetBillboardDailyFilterPageable(organizationId uuid.UUID, startDate int64, endDate int64, page int, size int, search string) ([]models.BillboardSummaryView, int64, error)
}

type billboardSummaryRepositoryImpl struct {
	db *gorm.DB
}

func NewBillBoardSummaryRepository() BillboardSummaryRepository {
	if db == nil {
		database.NewDatabaseConnection()
		db = database.GetDB()
	}
	return &billboardSummaryRepositoryImpl{
		db: db,
	}
}

func (r *billboardSummaryRepositoryImpl) GetStaffBillBoardsSummary(organizationId uuid.UUID, createdById uuid.UUID, page int, size int, search string) ([]models.BillboardSummaryView, int64, error) {
	var billboards []models.BillboardSummaryView
	var count int64
	err := r.db.Preload("Staff").Preload("Image").Preload("Campaign").Where("organization_id = ? AND created_by_id = ? AND board_code LIKE ?", organizationId, createdById, "%"+search+"%").Offset((page - 1) * size).Limit(size).Find(&billboards).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	return billboards, count, nil
}

func (r *billboardSummaryRepositoryImpl) GetStaffBillBoardsSummaryById(id uuid.UUID, createdById uuid.UUID) (*models.BillboardSummaryView, error) {
	var billboard models.BillboardSummaryView
	err := r.db.Preload("Staff").Preload("Image").Preload("Campaign").Where("billboard_id = ? AND created_by_id = ?", id, createdById).First(&billboard).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil

		}
		return nil, err
	}
	return &billboard, nil
}

func (r *billboardSummaryRepositoryImpl) GetBillboardDailyFilterPageable(organizationId uuid.UUID, startDate int64, endDate int64, page int, size int, search string) ([]models.BillboardSummaryView, int64, error) {
	var billboards []models.BillboardSummaryView
	var count int64
	err := r.db.Preload("Staff").Preload("Image").Preload("Campaign").Where("organization_id = ? AND created_at >= ? AND created_at <= ? AND board_code LIKE ?", organizationId, startDate, endDate, "%"+search+"%").Offset((page - 1) * size).Limit(size).Find(&billboards).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	return billboards, count, nil
}
