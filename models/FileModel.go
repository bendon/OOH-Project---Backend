package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FileModel struct {
	ID             uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	FileName       string     `gorm:"column:file_name;not null" json:"fileName"`
	FileUrl        string     `gorm:"column:file_url;not null" json:"fileUrl"`
	FileType       string     `gorm:"column:file_type;not null" json:"fileType"`
	FileSize       int64      `gorm:"column:file_size;not null" json:"fileSize"`
	FileExtension  string     `gorm:"column:file_extension;not null" json:"fileExtension"`
	OrganizationId *uuid.UUID `gorm:"type:char(36);column:organization_id;null" json:"organizationId"`
	UploadedById   *uuid.UUID `gorm:"type:char(36);column:uploaded_by_id; null" json:"uploadedById"`
	CreatedAt      int64      `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt      int64      `gorm:"column:updated_at" json:"updatedAt"`
}

func (FileModel) TableName() string {
	return "files"
}

func (f *FileModel) BeforeCreate(tx *gorm.DB) (err error) {
	f.ID = uuid.New()
	f.CreatedAt = time.Now().Unix()
	f.UpdatedAt = time.Now().Unix()
	return
}

func (f *FileModel) BeforeUpdate(tx *gorm.DB) (err error) {
	f.UpdatedAt = time.Now().Unix()
	return

}
