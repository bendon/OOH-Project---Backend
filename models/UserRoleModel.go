package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRoleModel struct {
	ID             uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	UserId         uuid.UUID `gorm:"type:char(36);column:user_id;not null" json:"userId"`
	RoleId         uuid.UUID `gorm:"type:char(36);column:role_Id" json:"roleId"`
	OrganizationId uuid.UUID `gorm:"type:char(36);column:organization_id" json:"organizationId"`
	CreatedAt      int64     `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt      int64     `gorm:"column:updated_at" json:"updatedAt"`
}

func (UserRoleModel) TableName() string {
	return "user_role"
}

func (userRole *UserRoleModel) BeforeCreate(tx *gorm.DB) (err error) {
	userRole.ID = uuid.New()
	userRole.CreatedAt = time.Now().Unix()
	userRole.UpdatedAt = time.Now().Unix()
	return
}

func (userRole *UserRoleModel) BeforeUpdate(tx *gorm.DB) (err error) {
	userRole.UpdatedAt = time.Now().Unix()
	return
}
