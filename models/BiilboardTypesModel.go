package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BillboardTypesModel struct {
	ID          uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	Name        string         `gorm:"type:varchar(255);column:name;not null" json:"name"`
	CreatedById uuid.UUID      `gorm:"type:char(36);column:created_by_id; null" json:"createdById"`
	DeletedById *uuid.UUID     `gorm:"type:char(36);column:created_by_id; null" json:"deletedById,omitempty"`
	CreatedAt   int64          `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt   int64          `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (BillboardTypesModel) TableName() string {
	return "bill_board_types"
}

func (b *BillboardTypesModel) BeforeCreate(tx *gorm.DB) (err error) {
	b.ID = uuid.New()
	b.CreatedAt = time.Now().Unix()
	b.UpdatedAt = time.Now().Unix()
	return
}

func (b *BillboardTypesModel) BeforeUpdate(tx *gorm.DB) (err error) {
	b.UpdatedAt = time.Now().Unix()
	return

}
