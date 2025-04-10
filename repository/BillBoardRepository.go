package repository

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"bbscout/config/database"
	"bbscout/models"
)

type BillBoardRepository interface {
	GetBillBoardById(id uuid.UUID) (*models.BillboardModel, error)
	CreateBillBoard(billboard *models.BillboardModel) (*models.BillboardModel, error)
	UpdateBillBoard(billboard *models.BillboardModel) (*models.BillboardModel, error)
	DeleteBillBoard(id uuid.UUID) error
	GetBillBoardByIdAndOrganizationId(id uuid.UUID, organizationId uuid.UUID) (*models.BillboardModel, error)
	GetBillBoardsByOrganizationId(organizationId uuid.UUID) ([]models.BillboardModel, error)
	GetBillBoardsByOrganizationIdPageable(organizationId uuid.UUID, page int, size int, search string) ([]models.BillboardModel, int64, error)
	GetBillboardUploadsByOrganizationIdAndCreatedById(organizationId uuid.UUID, createdById uuid.UUID, page int, size int, search string) ([]models.BillboardModel, int64, error)
	GetBillBoardBoundingBox(organizationId uuid.UUID, latitude float64, latitudeDef float64, longitude float64, longitudeDef float64) ([]models.BillboardModel, error)
}
type billBoardRepositoryImpl struct {
	db *gorm.DB
}

func NewBillBoardRepository() BillBoardRepository {
	if db == nil {
		database.NewDatabaseConnection()
		db = database.GetDB()
	}
	return &billBoardRepositoryImpl{
		db: db,
	}
}

func (r *billBoardRepositoryImpl) GetBillBoardById(id uuid.UUID) (*models.BillboardModel, error) {
	var billboard models.BillboardModel
	err := r.db.Where("id = ?", id).First(&billboard).Error
	if err != nil {
		return nil, err
	}
	return &billboard, nil
}

func (r *billBoardRepositoryImpl) CreateBillBoard(billboard *models.BillboardModel) (*models.BillboardModel, error) {
	err := r.db.Create(billboard).Error
	if err != nil {
		return nil, err
	}
	return billboard, nil
}

func (r *billBoardRepositoryImpl) UpdateBillBoard(billboard *models.BillboardModel) (*models.BillboardModel, error) {
	err := r.db.Save(billboard).Error
	if err != nil {
		return nil, err
	}
	return billboard, nil
}

func (r *billBoardRepositoryImpl) DeleteBillBoard(id uuid.UUID) error {
	err := r.db.Where("id = ?", id).Delete(&models.BillboardModel{}).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *billBoardRepositoryImpl) GetBillBoardByIdAndOrganizationId(id uuid.UUID, organizationId uuid.UUID) (*models.BillboardModel, error) {
	var billboard models.BillboardModel
	err := r.db.Preload("Image").Preload("CloseUpImage").Where("id = ? AND organization_id = ?", id, organizationId).First(&billboard).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil

		}
		return nil, err
	}
	return &billboard, nil
}

func (r *billBoardRepositoryImpl) GetBillBoardsByOrganizationId(organizationId uuid.UUID) ([]models.BillboardModel, error) {
	var billboards []models.BillboardModel
	err := r.db.Preload("Image").Where("organization_id = ?", organizationId).Find(&billboards).Error
	if err != nil {
		return nil, err
	}
	return billboards, nil
}

func (r *billBoardRepositoryImpl) GetBillBoardsByOrganizationIdPageable(organizationId uuid.UUID, page int, size int, search string) ([]models.BillboardModel, int64, error) {
	var billboards []models.BillboardModel
	var count int64
	err := r.db.Preload("Image").Where("organization_id = ? AND board_code LIKE ?", organizationId, "%"+search+"%").Offset((page - 1) * size).Limit(size).Find(&billboards).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	return billboards, count, nil
}

func (r *billBoardRepositoryImpl) GetBillboardUploadsByOrganizationIdAndCreatedById(organizationId uuid.UUID, createdById uuid.UUID, page int, size int, search string) ([]models.BillboardModel, int64, error) {
	var billboards []models.BillboardModel
	var count int64
	err := r.db.Preload("Image").Where("organization_id = ? AND created_by_id = ? AND board_code LIKE ?", organizationId, createdById, "%"+search+"%").Offset((page - 1) * size).Limit(size).Find(&billboards).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	return billboards, count, nil
}

func (r *billBoardRepositoryImpl) GetBillBoardBoundingBox(organizationId uuid.UUID, latitude float64, latitudeDef float64, longitude float64, longitudeDef float64) ([]models.BillboardModel, error) {
	var billboards []models.BillboardModel
	err := r.db.Preload("Image").Preload("CloseUpImage").Where("organization_id = ? AND latitude BETWEEN ? AND ? AND longitude BETWEEN ? AND ?", organizationId, latitude, latitudeDef, longitude, longitudeDef).Find(&billboards).Error
	if err != nil {
		return nil, err
	}
	return billboards, nil
}
