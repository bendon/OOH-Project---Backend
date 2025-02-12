package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PermissionModel struct {
	ID        uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(255);column:name;not null" json:"name"`
	Type      string    `gorm:"type:varchar(255);column:type;not null" json:"type"`
	Account   string    `gorm:"type:varchar(255);column:account;not null" json:"account"`
	CreatedAt int64     `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt int64     `gorm:"column:updated_at" json:"updatedAt"`
}

func (PermissionModel) TableName() string {
	return "permissions"
}

func (perm *PermissionModel) BeforeCreate(tx *gorm.DB) (err error) {
	perm.ID = uuid.New()
	perm.CreatedAt = time.Now().Unix()
	perm.UpdatedAt = time.Now().Unix()
	return
}

func (perm *PermissionModel) BeforeUpdate(tx *gorm.DB) (err error) {
	perm.UpdatedAt = time.Now().Unix()
	return
}
