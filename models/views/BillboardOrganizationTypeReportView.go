package models

import "github.com/google/uuid"

type BillboardOrganizationTypeReport struct {
	OrganizationID   uuid.UUID `gorm:"column:organization_id" json:"organizationId"`
	OrganizationName string    `gorm:"column:organization_name" json:"organizationName"`
	Type             string    `gorm:"column:type" json:"type"`
	TypeCount        int64     `gorm:"column:type_count" json:"typeCount"`
}

// TableName overrides the default table name for the GORM model
func (BillboardOrganizationTypeReport) TableName() string {
	return "billboard_organization_type_report"
}
