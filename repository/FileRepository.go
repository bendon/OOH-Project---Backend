package repository

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"bbscout/config/database"
	"bbscout/models"
)

type FileRepository interface {
	CreateFile(file *models.FileModel) (*models.FileModel, error)
	GetFileById(id uuid.UUID) (*models.FileModel, error)
	UpdateFile(file *models.FileModel) (*models.FileModel, error)
	DeleteFile(id uuid.UUID) error
	GetFiles() ([]models.FileModel, error)
}

type fileRepositoryImpl struct {
	db *gorm.DB
}

func NewFileRepository() FileRepository {
	if db == nil {
		database.NewDatabaseConnection()
		db = database.GetDB()
	}
	return &fileRepositoryImpl{
		db: db,
	}
}

func (r *fileRepositoryImpl) CreateFile(file *models.FileModel) (*models.FileModel, error) {
	err := r.db.Create(file).Error
	if err != nil {
		return nil, err
	}
	return file, nil
}
func (r *fileRepositoryImpl) GetFileById(id uuid.UUID) (*models.FileModel, error) {
	var file models.FileModel
	err := r.db.Where("id = ?", id).First(&file).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &file, nil
}
func (r *fileRepositoryImpl) UpdateFile(file *models.FileModel) (*models.FileModel, error) {
	err := r.db.Save(file).Error
	if err != nil {
		return nil, err
	}
	return file, nil
}
func (r *fileRepositoryImpl) DeleteFile(id uuid.UUID) error {
	err := r.db.Where("id = ?", id).Delete(&models.FileModel{}).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("file not found")
		}
		return err
	}
	return nil
}
func (r *fileRepositoryImpl) GetFiles() ([]models.FileModel, error) {
	var files []models.FileModel
	err := r.db.Find(&files).Error
	if err != nil {
		return nil, err
	}
	return files, nil
}
