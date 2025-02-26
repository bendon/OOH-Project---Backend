package models

import "github.com/google/uuid"

type BillboardUploadDayOfWeek struct {
	OrganizationID     uuid.UUID `gorm:"type:char(36);column:organization_id" json:"organizationId"`
	OrganizationName   string    `gorm:"type:varchar(255);column:organization_name" json:"organizationName"`
	OrganizationActive bool      `gorm:"type:boolean;column:organization_active" json:"organizationActive"`
	Year               int       `gorm:"type:int;column:upload_year" json:"year"`
	Month              int       `gorm:"type:int;column:upload_month" json:"month"`
	WeekNumber         int       `gorm:"type:int;column:upload_week_number" json:"weekNumber"`
	DayName            string    `gorm:"type:varchar(20);column:upload_day_name" json:"dayName"`
	TotalUploads       int64     `gorm:"type:bigint;column:total_uploads" json:"totalUploads"`
}

// TableName sets the table name for the view
func (BillboardUploadDayOfWeek) TableName() string {
	return "billboard_upload_day_of_week"
}
