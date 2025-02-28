package models

import "github.com/google/uuid"

type UserBillboardUploadMonthReport struct {
	UserID           uuid.UUID `gorm:"column:user_id" json:"userId"`
	UserName         string    `gorm:"column:user_name" json:"userName"`
	Email            string    `gorm:"column:email" json:"email"`
	OrganizationID   uuid.UUID `gorm:"column:organization_id" json:"organizationId"`
	OrganizationName string    `gorm:"column:organization_name" json:"organizationName"`
	UploadYear       int       `gorm:"column:upload_year" json:"uploadYear"`
	UploadMonth      int       `gorm:"column:upload_month" json:"uploadMonth"`
	BillboardCount   int       `gorm:"column:billboard_count" json:"billboardCount"`
}

// TableName overrides the table name for GORM
func (UserBillboardUploadMonthReport) TableName() string {
	return "user_billboard_uploads_by_month_year"
}
