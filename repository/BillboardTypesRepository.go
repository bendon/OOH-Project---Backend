package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"bbscout/config/database"
	"bbscout/models"
)

type BillboardTypesRepository interface {
	CreateBillboardType(billboardType *models.BillboardTypesModel) (*models.BillboardTypesModel, error)
	GetBillboardTypeByID(id uuid.UUID) (*models.BillboardTypesModel, error)
	GetBillboardTypes() ([]models.BillboardTypesModel, error)
	UpdateBillboardType(billboardType *models.BillboardTypesModel) error
	DeleteBillboardType(id uuid.UUID) error
	ExistsBillboardTypeByName(name string) (bool, error)
	ExistsBillboardTypeById(id uuid.UUID) (bool, error)
}

type BillboardTypesRepositoryImpl struct {
	DB *gorm.DB
}

func NewBillboardTypesRepository() BillboardTypesRepository {
	if db == nil {
		database.NewDatabaseConnection()
		db = database.GetDB()
	}
	return &BillboardTypesRepositoryImpl{
		DB: db,
	}
}

func (r *BillboardTypesRepositoryImpl) GetBillboardTypeByID(id uuid.UUID) (*models.BillboardTypesModel, error) {
	var billboardType models.BillboardTypesModel
	if err := r.DB.Where("id = ?", id).First(&billboardType).Error; err != nil {
		return nil, err
	}
	return &billboardType, nil
}
func (r *BillboardTypesRepositoryImpl) GetBillboardTypes() ([]models.BillboardTypesModel, error) {
	var billboardTypes []models.BillboardTypesModel
	if err := r.DB.Find(&billboardTypes).Error; err != nil {
		return nil, err
	}
	return billboardTypes, nil
}
func (r *BillboardTypesRepositoryImpl) UpdateBillboardType(billboardType *models.BillboardTypesModel) error {
	return r.DB.Save(billboardType).Error
}
func (r *BillboardTypesRepositoryImpl) DeleteBillboardType(id uuid.UUID) error {
	return r.DB.Delete(&models.BillboardTypesModel{}, id).Error
}
func (r *BillboardTypesRepositoryImpl) ExistsBillboardTypeByName(name string) (bool, error) {
	var count int64
	if err := r.DB.Model(&models.BillboardTypesModel{}).Where("name = ?", name).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *BillboardTypesRepositoryImpl) CreateBillboardType(billboardType *models.BillboardTypesModel) (*models.BillboardTypesModel, error) {
	if err := r.DB.Create(billboardType).Error; err != nil {
		return nil, err
	}
	return billboardType, nil
}
func (r *BillboardTypesRepositoryImpl) ExistsBillboardTypeById(id uuid.UUID) (bool, error) {
	var count int64
	if err := r.DB.Model(&models.BillboardTypesModel{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
