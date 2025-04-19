package models

import "github.com/google/uuid"

type UserMonthlyAuditStat struct {
	MonthNumber     int       `json:"monthNumber" gorm:"column:month_number"`
	MonthName       string    `json:"monthName" gorm:"column:month_name"`
	Year            int       `json:"year" gorm:"column:year_number"`
	UserId          uuid.UUID `json:"userId" gorm:"column:user_id"`
	OrganizationId  uuid.UUID `json:"organizationId" gorm:"column:organization_id"`
	TotalBillboards int64     `json:"totalBillboards" gorm:"column:total_billboards"`
}

// TableName overrides the default table name used by GORM
func (UserMonthlyAuditStat) TableName() string {
	return "user_monthly_audit_stats"
}
