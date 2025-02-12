package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrganizationUserModel struct {
	ID             uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	OrganizationId uuid.UUID  `gorm:"type:char(36);column:organization_id;not null" json:"organizationId"`
	UserId         uuid.UUID  `gorm:"type:char(36);column:user_id;not null" json:"userId"`
	CreatedById    *uuid.UUID `gorm:"type:char(36);column:created_by_id;null" json:"createdById"`
	Staff          *UserModel `gorm:"foreignKey:UserId; references:ID" json:"staff,omitempty"`
	CreatedAt      int64      `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt      int64      `gorm:"column:updated_at" json:"updatedAt"`
}

func (OrganizationUserModel) TableName() string {
	return "organization_user"
}

func (orgUser *OrganizationUserModel) BeforeCreate(tx *gorm.DB) (err error) {
	orgUser.ID = uuid.New()
	orgUser.CreatedAt = time.Now().Unix()
	orgUser.UpdatedAt = time.Now().Unix()
	return
}
func (orgUser *OrganizationUserModel) BeforeUpdate(tx *gorm.DB) (err error) {
	orgUser.UpdatedAt = time.Now().Unix()
	return
}
