package types

import (
	"github.com/google/uuid"
)

type CreateMonitoringRequest struct {
	BillboardId          *uuid.UUID        `json:"billboardId"` // billboard id
	County               *string           `json:"county"`
	Street               *string           `json:"street" `
	Location             *string           `json:"location" `
	Building             *string           `json:"building"`
	OwnerContacts        *Int64ArrayJSONB  `json:"ownerContacts"`
	OwnerEmail           *StringArrayJSONB `json:"ownerEmail"`
	Brand                *string           `json:"brand"`
	Campain              *string           `json:"campain"`
	Width                *float64          `json:"width" `
	Height               *float64          `json:"height" `
	Unit                 *string           `json:"unit" validate:"required,oneof=centimeters meters feet inches"`
	Structure            *string           `json:"structure"`
	Material             *string           `json:"material"`
	Angle                *string           `json:"angle"`
	Visibility           *string           `json:"visibility"`
	Illumination         *string           `json:"illumination"`
	ConditionOfMaterial  *string           `json:"conditionOfMaterial"`
	ConditionOfStructure *string           `json:"conditionOfStructure"`
	Comments             *string           `json:"comments"`
	Latitude             *float64          `json:"latitude"`
	Longitude            *float64          `json:"longitude"`
	Accuracy             *float64          `json:"accuracy" `
	CloseUpImageId       *uuid.UUID        `json:"closeUpImageId"`
	LongShotImageId      *uuid.UUID        `json:"longShotImageId"`
	Type                 *string           `json:"type"`
	Environment          *string           `json:"environment"`
}

type MonthlyMonitorReport struct {
	Month        int    `json:"month"`
	MonthName    string `json:"monthName"`
	TotalMonitor int64  `json:"totalMonitor"`
}

type UserReportResponse struct {
	Year     int                    `json:"year"`
	Monitors MonitorsReportResponse `json:"monitor"`
	Audit    AuditReportResponse    `json:"audit"`
}

type MonitorsReportResponse struct {
	Monitor       interface{}            `json:"summary"`
	MonthlyReport []MonthlyMonitorReport `json:"monthlyReport"`
}
