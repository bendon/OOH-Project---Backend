package models

import (
	"github.com/google/uuid"

	"bbscout/models"
)

type UserBillboardUploadReport struct {
	UserId           uuid.UUID         `gorm:"column:user_id" json:"userId"`
	UserName         string            `gorm:"column:user_name" json:"userName"`
	Email            string            `gorm:"column:email" json:"email"`
	OrganizationID   uuid.UUID         `gorm:"column:organization_id" json:"organizationId"`
	OrganizationName string            `gorm:"column:organization_name" json:"organizationName"`
	BillboardCount   int               `gorm:"column:billboard_count" json:"billboardCount"`
	UserInfo         *models.UserModel `gorm:"foreignKey:UserId; references:ID" json:"userInfo,omitempty"`
}

// TableName sets the name of the view in the database
func (UserBillboardUploadReport) TableName() string {
	return "user_billboard_uploads_report"
}
