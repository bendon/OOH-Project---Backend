package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type CreateBillboardRequest struct {
	Description     string            `json:"description"`
	ParentBoardCode *string           `json:"parentBoardCode"`
	ImageID         uuid.UUID         `json:"imageId"`
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
}

type UpdateBillboardRequest struct {
	Description     string            `json:"description"`
	ParentBoardCode *string           `json:"parentBoardCode"`
	ImageID         uuid.UUID         `json:"imageId"`
	Location        string            `json:"location" validate:"required"`
	Latitude        float64           `json:"latitude" validate:"required,latitude"`
	Longitude       float64           `json:"longitude" validate:"required,longitude"`
	Accuracy        float64           `json:"accuracy" validate:"required"`
	Width           float64           `json:"width" validate:"required,gt=0"`
	Height          float64           `json:"height" validate:"required,gt=0"`
	Unit            string            `json:"unit" validate:"required,oneof=centimeters meters feet inches"`
	Type            string            `json:"type" validate:"required,oneof=digital static LED traditional"`
	Price           float64           `json:"price" validate:"required,gt=0"`
	ObjectType      *string           `json:"objectType"`
	Owner           *string           `json:"owner"`
	OwnerContacts   *Int64ArrayJSONB  `json:"ownerContacts"`
	OwnerEmail      *StringArrayJSONB `json:"ownerEmail"`
	Occupied        bool              `json:"occupied" default:"false"`
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
