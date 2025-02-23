package repository

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"bbscout/config/database"
	"bbscout/models"
)

type BillboardCampaignRepository interface {
	CreateBillboardCampaign(campaign *models.BillboardCampaignModel) (*models.BillboardCampaignModel, error)
	GetBillboardCampaignById(id uuid.UUID) (*models.BillboardCampaignModel, error)
	GetBillboardCampaignByIdAndOrganizationId(id uuid.UUID, organizationId uuid.UUID) (*models.BillboardCampaignModel, error)
	GetBillboardCampaignByBillboardId(billboardId uuid.UUID) ([]models.BillboardCampaignModel, error)
	GetBillboardCampaignByBillboardIdAndOrganizationId(billboardId uuid.UUID, organizationId uuid.UUID) ([]models.BillboardCampaignModel, error)
	GetBillboardCampaignsByOrganizationId(organizationId uuid.UUID) ([]models.BillboardCampaignModel, error)
	UpdateBillboardCampaign(campaign *models.BillboardCampaignModel) (*models.BillboardCampaignModel, error)
	DeleteBillboardCampaign(id uuid.UUID) error
	FindBillboardCampaignByOrganizationId(organizationId uuid.UUID, page int, size int, search string) ([]models.BillboardCampaignModel, int64, error)
	FindBillboardCampaignByOrganizationIdAndBillboardIdPageable(organizationId uuid.UUID, billboardId uuid.UUID, page int, size int, search string) ([]models.BillboardCampaignModel, int64, error)
	FindBillboardCampaignByOrganizationIdAndBillboardIdAndActive(organizationId uuid.UUID, billboardId uuid.UUID, active bool) (*models.BillboardCampaignModel, error)
}
type BillboardCampaignRepositoryImpl struct {
	db *gorm.DB
}

func NewBillboardCampaignRepository() BillboardCampaignRepository {
	if db == nil {
		database.NewDatabaseConnection()
		db = database.GetDB()
	}
	return &BillboardCampaignRepositoryImpl{
		db: db,
	}
}

func (r *BillboardCampaignRepositoryImpl) CreateBillboardCampaign(campaign *models.BillboardCampaignModel) (*models.BillboardCampaignModel, error) {
	err := r.db.Create(campaign).Error
	if err != nil {
		return nil, err
	}
	return campaign, nil
}
func (r *BillboardCampaignRepositoryImpl) GetBillboardCampaignById(id uuid.UUID) (*models.BillboardCampaignModel, error) {
	var campaign models.BillboardCampaignModel
	err := r.db.Where("id = ?", id).First(&campaign).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil

		}
		return nil, err
	}
	return &campaign, nil
}
func (r *BillboardCampaignRepositoryImpl) GetBillboardCampaignByBillboardId(billboardId uuid.UUID) ([]models.BillboardCampaignModel, error) {
	var campaign []models.BillboardCampaignModel
	err := r.db.Where("billboard_id = ?", billboardId).Find(&campaign).Error
	if err != nil {
		return nil, err
	}
	return campaign, nil
}
func (r *BillboardCampaignRepositoryImpl) GetBillboardCampaignByBillboardIdAndOrganizationId(billboardId uuid.UUID, organizationId uuid.UUID) ([]models.BillboardCampaignModel, error) {
	var campaign []models.BillboardCampaignModel
	err := r.db.Where("billboard_id = ? AND organization_id = ?", billboardId, organizationId).Find(&campaign).Error
	if err != nil {
		return nil, err
	}
	return campaign, nil
}
func (r *BillboardCampaignRepositoryImpl) GetBillboardCampaignsByOrganizationId(organizationId uuid.UUID) ([]models.BillboardCampaignModel, error) {
	var campaigns []models.BillboardCampaignModel
	err := r.db.Where("organization_id = ?", organizationId).Find(&campaigns).Error
	if err != nil {
		return nil, err
	}
	return campaigns, nil
}
func (r *BillboardCampaignRepositoryImpl) UpdateBillboardCampaign(campaign *models.BillboardCampaignModel) (*models.BillboardCampaignModel, error) {
	err := r.db.Save(campaign).Error
	if err != nil {
		return nil, err
	}
	return campaign, nil
}
func (r *BillboardCampaignRepositoryImpl) DeleteBillboardCampaign(id uuid.UUID) error {
	err := r.db.Where("id = ?", id).Delete(&models.BillboardCampaignModel{}).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *BillboardCampaignRepositoryImpl) FindBillboardCampaignByOrganizationId(organizationId uuid.UUID, page int, size int, search string) ([]models.BillboardCampaignModel, int64, error) {
	var campaigns []models.BillboardCampaignModel
	var total int64
	err := r.db.Preload("Billboard").Preload("Image").Where("organization_id = ?", organizationId).Offset((page - 1) * size).Limit(size).Find(&campaigns).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	return campaigns, total, nil
}
func (r *BillboardCampaignRepositoryImpl) FindBillboardCampaignByOrganizationIdAndBillboardIdPageable(organizationId uuid.UUID, billboardId uuid.UUID, page int, size int, search string) ([]models.BillboardCampaignModel, int64, error) {
	var campaigns []models.BillboardCampaignModel
	var total int64
	err := r.db.Where("organization_id = ? AND billboard_id = ?", organizationId, billboardId).Offset((page - 1) * size).Limit(size).Find(&campaigns).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	return campaigns, total, nil
}
func (r *BillboardCampaignRepositoryImpl) FindBillboardCampaignByOrganizationIdAndBillboardIdAndActive(organizationId uuid.UUID, billboardId uuid.UUID, active bool) (*models.BillboardCampaignModel, error) {
	var campaign models.BillboardCampaignModel
	err := r.db.Where("organization_id = ? AND billboard_id = ? AND active = ?", organizationId, billboardId, active).First(&campaign).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &campaign, nil
}

func (r *BillboardCampaignRepositoryImpl) GetBillboardCampaignByIdAndOrganizationId(id uuid.UUID, organizationId uuid.UUID) (*models.BillboardCampaignModel, error) {
	var campaign models.BillboardCampaignModel
	err := r.db.Preload("Billboard").Preload("Image").Where("id = ? AND organization_id = ?", id, organizationId).First(&campaign).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &campaign, nil
}
