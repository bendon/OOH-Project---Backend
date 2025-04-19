package models

import "github.com/google/uuid"

type UserMonthlyMonitoringStat struct {
	MonthNumber    int       `json:"monthNumber" gorm:"column:month_number"`
	MonthName      string    `json:"monthName" gorm:"column:month_name"`
	Year           int       `json:"year" gorm:"column:year_number"`
	UserId         uuid.UUID `json:"userId" gorm:"column:user_id"`
	OrganizationId uuid.UUID `json:"organizationId" gorm:"column:organization_id"`
	TotalMonitors  int64     `json:"totalMonitors" gorm:"column:total_monitors"`
}

// TableName overrides the default table name used by GORM
func (UserMonthlyMonitoringStat) TableName() string {
	return "user_monthly_monitoring_stats"
}
