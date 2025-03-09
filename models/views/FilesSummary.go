package models

import (
	"github.com/google/uuid"

	"bbscout/models"
)

type FilesSummary struct {
	FileID         uuid.UUID              `gorm:"column:file_id" json:"file_id"`
	FileName       string                 `gorm:"column:file_name" json:"file_name"`
	FileExtension  string                 `gorm:"column:file_extension" json:"file_extension"`
	FileType       string                 `gorm:"column:file_type" json:"file_type"`
	FileURL        string                 `gorm:"column:file_url" json:"file_url"`
	OrganizationID uuid.UUID              `gorm:"column:organization_id" json:"organization_id"`
	FileSize       int64                  `gorm:"column:file_size" json:"file_size"`
	CreatedAt      int64                  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      int64                  `gorm:"column:updated_at" json:"updated_at"`
	BillboardID    *uuid.UUID             `gorm:"column:billboard_id" json:"billboard_id"` // Use *uint to handle nullable BillboardID
	Billboard      *models.BillboardModel `gorm:"foreignKey:BillboardID; references:ID" json:"billboard,omitempty"`
}

// TableName overrides the table name used by GORM
func (FilesSummary) TableName() string {
	return "files_summary"
}
