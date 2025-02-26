package models

import "github.com/google/uuid"

type BillboardUploadMonthlyReport struct {
	OrganizationID          uuid.UUID `gorm:"type:char(36);column:organization_id" json:"organizationId"`
	OrganizationName        string    `gorm:"type:varchar(255);column:organization_name" json:"organizationName"`
	OrganizationDescription string    `gorm:"type:varchar(2000);column:organization_description" json:"organizationDescription"`
	OrganizationActive      bool      `gorm:"type:boolean;column:organization_active" json:"organizationActive"`
	Year                    int       `gorm:"type:int;column:upload_year" json:"year"`
	Month                   int       `gorm:"type:int;column:upload_month" json:"month"`
	MonthName               string    `gorm:"type:varchar(20);column:upload_month_name" json:"monthName"`
	TotalUploads            int64     `gorm:"type:bigint;column:total_uploads" json:"totalUploads"`
}

// TableName sets the table name for the view
func (BillboardUploadMonthlyReport) TableName() string {
	return "billboard_upload_monthly_report"
}
