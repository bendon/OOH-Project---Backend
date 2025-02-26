package models

import "github.com/google/uuid"

type BillboardOrganizationLocationReport struct {
	OrganizationID   uuid.UUID `gorm:"column:organization_id" json:"organizationId"`
	OrganizationName string    `gorm:"column:organization_name" json:"organizationName"`
	Location         string    `gorm:"column:location" json:"location"`
	CountPerLocation int64     `gorm:"column:count_per_location" json:"countPerLocation"`
}

// TableName overrides the default table name for the GORM model
func (BillboardOrganizationLocationReport) TableName() string {
	return "billboard_organization_location_report"
}
