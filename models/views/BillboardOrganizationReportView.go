package models

import "github.com/google/uuid"

type BillboardOrganizationReport struct {
	OrganizationID          uuid.UUID `gorm:"column:organization_id" json:"organizationId"`
	OrganizationName        string    `gorm:"column:organization_name" json:"organizationName"`
	OrganizationDescription string    `gorm:"column:organization_description" json:"organizationDescription"`
	TotalUploads            int64     `gorm:"column:total_uploads" json:"totalUploads"`
	Occupied                int64     `gorm:"column:total_occupied" json:"totalOccupied"`
	NotOccupied             int64     `gorm:"column:total_not_occupied" json:"notOccupied"`
	Today                   int64     `gorm:"column:uploaded_today" json:"today"`
	ThisMonth               int64     `gorm:"column:uploaded_this_month" json:"thisMonth"`
	ThisYear                int64     `gorm:"column:uploaded_this_year" json:"thisYear"`
}

// TableName overrides the default table name for the GORM model
func (BillboardOrganizationReport) TableName() string {
	return "billboard_organization_report"
}
