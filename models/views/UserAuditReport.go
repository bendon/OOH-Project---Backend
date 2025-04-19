package models

import "github.com/google/uuid"

type UserAuditReport struct {
	UserId          uuid.UUID `json:"userId" gorm:"column:user_id"`
	OrganizationId  uuid.UUID `json:"organizationId" gorm:"column:organization_id"`
	TotalBillboards int64     `json:"totalBillboards" gorm:"column:total_billboards"`
}

// TableName overrides the default table name used by GORM
func (UserAuditReport) TableName() string {
	return "user_audit_report"
}
