package repository

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"bbscout/config/database"
	"bbscout/models"
)

type BillBoardSequenceRepository interface {
	GetBillBoardSequenceByOrganizationId(organizationId uuid.UUID) (*models.BillboardSequenceModel, error)
	CreateBillBoardSequence(billboardSequence *models.BillboardSequenceModel) (*models.BillboardSequenceModel, error)
	UpdateBillBoardSequence(billboardSequence *models.BillboardSequenceModel) (*models.BillboardSequenceModel, error)
	GetBillBoardByIdAndOrganizationId(id uuid.UUID, organizationId uuid.UUID) (*models.BillboardSequenceModel, error)
}
type billBoardSequenceRepositoryImpl struct {
	db *gorm.DB
}

func NewBillBoardSequenceRepository() BillBoardSequenceRepository {
	if db == nil {
		database.NewDatabaseConnection()
		db = database.GetDB()
	}
	return &billBoardSequenceRepositoryImpl{
		db: db,
	}
}
func (r *billBoardSequenceRepositoryImpl) GetBillBoardSequenceByOrganizationId(organizationId uuid.UUID) (*models.BillboardSequenceModel, error) {
	var billboardSequence models.BillboardSequenceModel
	err := r.db.Where("organization_id = ?", organizationId).First(&billboardSequence).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil

		}
		return nil, err
	}
	return &billboardSequence, nil
}
func (r *billBoardSequenceRepositoryImpl) CreateBillBoardSequence(billboardSequence *models.BillboardSequenceModel) (*models.BillboardSequenceModel, error) {
	err := r.db.Create(billboardSequence).Error
	if err != nil {
		return nil, err
	}
	return billboardSequence, nil
}
func (r *billBoardSequenceRepositoryImpl) UpdateBillBoardSequence(billboardSequence *models.BillboardSequenceModel) (*models.BillboardSequenceModel, error) {
	err := r.db.Save(billboardSequence).Error
	if err != nil {
		return nil, err
	}
	return billboardSequence, nil
}
func (r *billBoardSequenceRepositoryImpl) GetBillBoardByIdAndOrganizationId(id uuid.UUID, organizationId uuid.UUID) (*models.BillboardSequenceModel, error) {
	var billboardSequence models.BillboardSequenceModel
	err := r.db.Where("id = ? AND organization_id = ?", id, organizationId).First(&billboardSequence).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil

		}
		return nil, err
	}
	return &billboardSequence, nil
}
