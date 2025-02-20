package repository

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"bbscout/config/database"
	models "bbscout/models/views"
)

type OrganizationUserSummaryRepository interface {
	GetOrganizationUserSummary(organizationID uuid.UUID) ([]models.OrganizationUserSummary, error)
	GetOrganizationUserSummaryPageable(organizationID uuid.UUID, page int, size int, search string) ([]models.OrganizationUserSummary, int64, error)
	GetOrganizationUserSummaryByUserId(userId uuid.UUID, organizationID uuid.UUID) (*models.OrganizationUserSummary, error)
}

type organizationUserSummaryRepositoryImpl struct {
	db *gorm.DB
}

func NewOrganizationUserSummaryRepository() OrganizationUserSummaryRepository {
	if db == nil {
		database.NewDatabaseConnection()
		db = database.GetDB()
	}
	return &organizationUserSummaryRepositoryImpl{
		db: db,
	}
}

func (r *organizationUserSummaryRepositoryImpl) GetOrganizationUserSummary(organizationID uuid.UUID) ([]models.OrganizationUserSummary, error) {
	var organizationUserSummary []models.OrganizationUserSummary
	err := r.db.Where("organization_id = ?", organizationID).Find(&organizationUserSummary).Error
	if err != nil {
		return nil, err
	}
	return organizationUserSummary, nil
}
func (r *organizationUserSummaryRepositoryImpl) GetOrganizationUserSummaryByUserId(userId uuid.UUID, organizationID uuid.UUID) (*models.OrganizationUserSummary, error) {
	var organizationUserSummary models.OrganizationUserSummary
	err := r.db.Where("user_id = ? AND organization_id = ?", userId, organizationID).First(&organizationUserSummary).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil

		}
		return nil, err
	}
	return &organizationUserSummary, nil
}
func (r *organizationUserSummaryRepositoryImpl) GetOrganizationUserSummaryPageable(organizationId uuid.UUID, page int, size int, search string) ([]models.OrganizationUserSummary, int64, error) {
	var organizationUserSummary []models.OrganizationUserSummary
	var totalCount int64

	if err := r.db.Where("organization_id = ? ", organizationId).Model(&models.OrganizationUserSummary{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * size

	if err := r.db.Where("organization_id = ?  AND first_name LIKE ? OR last_name LIKE ?", organizationId, "%"+search+"%", "%"+search+"%").
		Limit(size).
		Offset(offset).
		Find(&organizationUserSummary).Error; err != nil {
		return nil, 0, err
	}
	return organizationUserSummary, totalCount, nil

}
