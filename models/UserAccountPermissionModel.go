package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserAccountPermissionModel struct {
	ID             uuid.UUID        `gorm:"type:char(36);primaryKey" json:"id"`
	AccountId      uuid.UUID        `gorm:"type:char(36);column:account_id;not null" json:"accountId"`
	OrganizationId uuid.UUID        `gorm:"type:char(36);column:organization_id;not null" json:"orgnizationId"`
	UserId         uuid.UUID        `gorm:"type:char(36);column:user_id;not null" json:"userId"`
	PermissionId   uuid.UUID        `gorm:"type:char(36);column:permission_id;not null" json:"permissionId"`
	Permission     *PermissionModel `gorm:"foreignKey:PermissionId; references:ID" json:"permission,omitempty"`
	CreatedAt      int64            `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt      int64            `gorm:"column:updated_at" json:"updatedAt"`
}

func (UserAccountPermissionModel) TableName() string {
	return "user_account_permissions"
}

func (accPermission *UserAccountPermissionModel) BeforeCreate(tx *gorm.DB) (err error) {
	accPermission.ID = uuid.New()
	accPermission.CreatedAt = time.Now().Unix()
	accPermission.UpdatedAt = time.Now().Unix()
	return
}
