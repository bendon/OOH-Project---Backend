package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math"

	"github.com/google/uuid"
)

type CreateBillboardRequest struct {
	Description     string            `json:"description"`
	ParentBoardCode *string           `json:"parentBoardCode"`
	ImageID         *uuid.UUID        `json:"imageId"`
	Location        string            `json:"location" validate:"required"`
	Latitude        float64           `json:"latitude" validate:"required,latitude"`
	Longitude       float64           `json:"longitude" validate:"required,longitude"`
	Accuracy        float64           `json:"accuracy" `
	Width           float64           `json:"width" `
	Height          float64           `json:"height" `
	Unit            string            `json:"unit" validate:"required,oneof=centimeters meters feet inches"`
	Type            string            `json:"type"`
	Price           float64           `json:"price" `
	ObjectType      *string           `json:"objectType"`
	Occupied        bool              `json:"occupied" default:"false"`
	Owner           *string           `json:"owner"`
	OwnerContacts   *Int64ArrayJSONB  `json:"ownerContacts"`
	OwnerEmail      *StringArrayJSONB `json:"ownerEmail"`
	City            *string           `json:"city"`
	CloseUpImageId  uuid.UUID         `json:"closeUpImageId"`
	Structure       *string           `json:"structure"`
	Material        *string           `json:"material"`
	Angle           *string           `json:"angle"`
	Visibility      *string           `json:"visibility"`
	Illumination    *string           `json:"illumination"`
}

type UpdateBillboardRequest struct {
	Description     *string           `json:"description"`
	ParentBoardCode *string           `json:"parentBoardCode"`
	ImageID         *uuid.UUID        `json:"imageId"`
	Location        *string           `json:"location"`
	Latitude        *float64          `json:"latitude"`
	Longitude       *float64          `json:"longitude"`
	Accuracy        *float64          `json:"accuracy"`
	Width           *float64          `json:"width"`
	Height          *float64          `json:"height"`
	Unit            string            `json:"unit" validate:"required,oneof=centimeters meters feet inches"`
	Type            string            `json:"type"`
	Price           *float64          `json:"price"`
	ObjectType      *string           `json:"objectType"`
	Owner           *string           `json:"owner"`
	OwnerContacts   *Int64ArrayJSONB  `json:"ownerContacts"`
	OwnerEmail      *StringArrayJSONB `json:"ownerEmail"`
	Occupied        bool              `json:"occupied" default:"false"`
	City            *string           `json:"city"`
	CloseUpImageId  uuid.UUID         `json:"closeUpImageId"`
	Structure       *string           `json:"structure"`
	Material        *string           `json:"material"`
	Angle           *string           `json:"angle"`
	Visibility      *string           `json:"visibility"`
	Illumination    *string           `json:"illumination"`
}

type BillboardTypeRequest struct {
	Name string `json:"name" validate:"required"`
}

type JSONB map[string]interface{}

// Marshal JSON before saving to DB
func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return "{}", nil
	}
	return json.Marshal(j)
}

// Unmarshal JSON when retrieving from DB
func (j *JSONB) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONB value: %v", value)
	}
	return json.Unmarshal(bytes, &j)
}

type StringArrayJSONB []string

// Marshal JSON before saving to DB
func (s StringArrayJSONB) Value() (driver.Value, error) {
	if s == nil {
		return "[]", nil
	}
	return json.Marshal(s)
}

// Unmarshal JSON when retrieving from DB
func (s *StringArrayJSONB) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONB string array: %v", value)
	}
	return json.Unmarshal(bytes, s)
}

type Int64ArrayJSONB []int64

// Marshal JSON before saving to DB
func (s Int64ArrayJSONB) Value() (driver.Value, error) {
	if s == nil {
		return "[]", nil
	}
	return json.Marshal(s)
}

// Unmarshal JSON when retrieving from DB
func (s *Int64ArrayJSONB) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONB int64 array: %v", value)
	}
	return json.Unmarshal(bytes, s)
}

type LocationRequest struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// calculateDistance calculates the distance between two points using the Haversine formula
// Returns distance in meters
func CalculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371000 // Earth's radius in meters

	// Convert degrees to radians
	lat1Rad := lat1 * math.Pi / 180
	lon1Rad := lon1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lon2Rad := lon2 * math.Pi / 180

	// Haversine formula
	dLat := lat2Rad - lat1Rad
	dLon := lon2Rad - lon1Rad
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	distance := earthRadius * c

	return distance
}

type MonthlyAuditReport struct {
	Month      int    `json:"month"`
	MonthName  string `json:"monthName"`
	TotalAudit int64  `json:"totalAudit"`
}

type AuditReportResponse struct {
	Audit         interface{}          `json:"summary"`
	MonthlyReport []MonthlyAuditReport `json:"monthlyReport"`
}
