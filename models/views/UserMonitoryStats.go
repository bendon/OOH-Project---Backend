package models

import "github.com/google/uuid"

type UserMonitoringStat struct {
	UserID         uuid.UUID `json:"userId" gorm:"column:user_id"`
	OrganizationId uuid.UUID `json:"organizationId" gorm:"column:organization_id"`
	TotalMonitors  int64     `json:"totalMonitors" gorm:"column:total_monitors"`
}

// TableName overrides the default table name used by GORM
func (UserMonitoringStat) TableName() string {
	return "user_monitoring_stats"
}
