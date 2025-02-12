package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrganizationModel struct {
	ID          uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	Name        string     `gorm:"type:varchar(255);column:name;not null" json:"name"`
	Description string     `gorm:"type:varchar(2000);column:description;not null" json:"description"`
	IsActive    bool       `gorm:"type:boolean;column:is_active;default:false;not null" json:"isActive"`
	AdminId     uuid.UUID  `gorm:"type:char(36);column:admin_id;not null" json:"-"`
	IsOperation bool       `gorm:"type:boolean;default:false;column:is_operation;not null" json:"-"`
	Admin       *UserModel `gorm:"foreignKey:AdminId; references:ID" json:"admin,omitempty"`
	CreatedAt   int64      `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt   int64      `gorm:"column:updated_at" json:"updatedAt"`
}

func (OrganizationModel) TableName() string {
	return "organization"
}

func (org *OrganizationModel) BeforeCreate(tx *gorm.DB) (err error) {
	org.ID = uuid.New()
	org.CreatedAt = time.Now().Unix()
	org.UpdatedAt = time.Now().Unix()
	return
}
func (org *OrganizationModel) BeforeUpdate(tx *gorm.DB) (err error) {
	org.UpdatedAt = time.Now().Unix()
	return
}
