package repository

import (
	"bbscout/config/database"
	"bbscout/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MonitoringRepository interface {
	CreateMonitoring(monitoring *models.MonitoringModel) (*models.MonitoringModel, error)
	GetMonitoring(id uuid.UUID) (*models.MonitoringModel, error)
	GetAllMonitoring(organizationId uuid.UUID, page, size int) ([]*models.MonitoringModel, int64, error)
	GetMonitoringByUser(organizationId uuid.UUID, userId uuid.UUID, page, size int) ([]*models.MonitoringModel, int64, error)
	UpdateMonitoring(monitoring *models.MonitoringModel) (*models.MonitoringModel, error)
	DeleteMonitoring(id uuid.UUID) error
}

type monitoringRepositoryImpl struct {
	db *gorm.DB
}

func NewMonitoringRepository() MonitoringRepository {
	if db == nil {
		database.NewDatabaseConnection()
		db = database.GetDB()
	}
	return &monitoringRepositoryImpl{
		db: db,
	}
}

func (m *monitoringRepositoryImpl) CreateMonitoring(monitoring *models.MonitoringModel) (*models.MonitoringModel, error) {
	if err := m.db.Create(monitoring).Error; err != nil {
		return nil, err
	}
	return monitoring, nil
}
func (m *monitoringRepositoryImpl) GetMonitoring(id uuid.UUID) (*models.MonitoringModel, error) {
	var monitoring models.MonitoringModel
	if err := m.db.First(&monitoring, id).Error; err != nil {
		return nil, err
	}
	return &monitoring, nil
}

func (m *monitoringRepositoryImpl) UpdateMonitoring(monitoring *models.MonitoringModel) (*models.MonitoringModel, error) {
	if err := m.db.Save(monitoring).Error; err != nil {
		return nil, err
	}
	return monitoring, nil
}
func (m *monitoringRepositoryImpl) DeleteMonitoring(id uuid.UUID) error {
	if err := m.db.Delete(&models.MonitoringModel{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (m *monitoringRepositoryImpl) GetAllMonitoring(organizationId uuid.UUID, page, size int) ([]*models.MonitoringModel, int64, error) {
	var monitoring []*models.MonitoringModel
	var totalCount int64
	if err := m.db.Where("organization_id = ?", organizationId).Model(&models.MonitoringModel{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * size
	if err := m.db.Preload("LongShortImage").Preload("CloseUpImage").Preload("User").Preload("Billboard").Where("organization_id = ?", organizationId).
		Order("created_at DESC").
		Limit(size).
		Offset(offset).
		Find(&monitoring).Error; err != nil {
		return nil, 0, err
	}
	return monitoring, totalCount, nil
}

func (m *monitoringRepositoryImpl) GetMonitoringByUser(organizationId uuid.UUID, userId uuid.UUID, page, size int) ([]*models.MonitoringModel, int64, error) {
	var monitoring []*models.MonitoringModel
	var totalCount int64
	if err := m.db.Where("organization_id = ? AND monitored_by_id = ?", organizationId, userId).Model(&models.MonitoringModel{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * size
	if err := m.db.Preload("LongShortImage").Preload("CloseUpImage").Preload("User").Preload("Billboard").Where("organization_id = ? AND monitored_by_id = ?", organizationId, userId).
		Order("created_at DESC").
		Limit(size).
		Offset(offset).
		Find(&monitoring).Error; err != nil {
		return nil, 0, err
	}
	return monitoring, totalCount, nil
}
