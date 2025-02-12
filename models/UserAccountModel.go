package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserAccountModel struct {
	ID             uuid.UUID          `gorm:"type:char(36);primaryKey" json:"id"`
	OrganizationId *uuid.UUID         `gorm:"type:char(36);column:organization_id;null" json:"-"`
	UserId         uuid.UUID          `gorm:"type:char(36);column:user_id;not null" json:"userId"`
	Active         bool               `gorm:"type:boolean;column:active;not null" json:"active"`
	IsLocked       bool               `gorm:"type:boolean;column:is_locked;not null" json:"isLocked"`
	Organization   *OrganizationModel `gorm:"foreignKey:OrganizationId; references:ID" json:"organization,omitempty"`
	Staff          *UserModel         `gorm:"foreignKey:UserId; references:ID" json:"staff,omitempty"`
	CreatedAt      int64              `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt      int64              `gorm:"column:updated_at" json:"updatedAt"`
}

func (UserAccountModel) TableName() string {
	return "user_account"
}

func (userAcc *UserAccountModel) BeforeCreate(tx *gorm.DB) (err error) {
	userAcc.ID = uuid.New()
	userAcc.CreatedAt = time.Now().Unix()
	userAcc.UpdatedAt = time.Now().Unix()
	return
}
func (userAcc *UserAccountModel) BeforeUpdate(tx *gorm.DB) (err error) {
	userAcc.UpdatedAt = time.Now().Unix()
	return
}
