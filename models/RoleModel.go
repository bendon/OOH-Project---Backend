package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoleModel struct {
	ID             uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	Name           string     `gorm:"column:name;not null" json:"name"`
	OrganizationId *uuid.UUID `gorm:"type:char(36);column:organization_id;null" json:"organizationId"`
	CreatedAt      int64      `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt      int64      `gorm:"column:updated_at" json:"updatedAt"`
}

func (RoleModel) TableName() string {
	return "roles"
}

func (role *RoleModel) BeforeCreate(tx *gorm.DB) (err error) {
	role.ID = uuid.New()
	role.CreatedAt = time.Now().Unix()
	role.UpdatedAt = time.Now().Unix()
	return
}

func (role *RoleModel) BeforeUpdate(tx *gorm.DB) (err error) {
	role.UpdatedAt = time.Now().Unix()
	return

}
