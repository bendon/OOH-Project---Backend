package models

import "github.com/google/uuid"

type UserBillboardUploadsReportWeek struct {
	UserID           uuid.UUID `gorm:"column:user_id" json:"userId"`
	UserName         string    `gorm:"column:user_name" json:"userName"`
	Email            string    `gorm:"column:email" json:"email"`
	OrganizationID   uuid.UUID `gorm:"column:organization_id" json:"organizationId"`
	OrganizationName string    `gorm:"column:organization_name" json:"organizationName"`
	BillboardCount   int       `gorm:"column:billboard_count" json:"billboardCount"`
	UploadYear       int       `gorm:"column:upload_year" json:"uploadYear"`
	UploadMonth      int       `gorm:"column:upload_month" json:"uploadMonth"`
	WeekNumber       int       `gorm:"column:week_number" json:"weekNumber"`
	DayOfWeek        string    `gorm:"column:day_of_week" json:"dayOfWeek"`
}

// TableName overrides the table name for GORM
func (UserBillboardUploadsReportWeek) TableName() string {
	return "user_billboard_uploads_report_week"
}
