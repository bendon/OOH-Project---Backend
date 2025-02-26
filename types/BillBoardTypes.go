package types

import "github.com/google/uuid"

type CreateBillboardRequest struct {
	Description     string    `json:"description"`
	ParentBoardCode *string   `json:"parentBoardCode"`
	ImageID         uuid.UUID `json:"imageId"`
	Location        string    `json:"location" validate:"required"`
	Latitude        float64   `json:"latitude" validate:"required,latitude"`
	Longitude       float64   `json:"longitude" validate:"required,longitude"`
	Accuracy        float64   `json:"accuracy" validate:"gte=0"`
	Width           float64   `json:"width" validate:"required,gt=0"`
	Height          float64   `json:"height" validate:"required,gt=0"`
	Unit            string    `json:"unit" validate:"required,oneof=centimeters meters feet inches"`
	Type            string    `json:"type"`
	Price           float64   `json:"price" `
}

type UpdateBillboardRequest struct {
	Description     string    `json:"description"`
	ParentBoardCode *string   `json:"parentBoardCode"`
	ImageID         uuid.UUID `json:"imageId"`
	Location        string    `json:"location" validate:"required"`
	Latitude        float64   `json:"latitude" validate:"required,latitude"`
	Longitude       float64   `json:"longitude" validate:"required,longitude"`
	Accuracy        float64   `json:"accuracy" validate:"required"`
	Width           float64   `json:"width" validate:"required,gt=0"`
	Height          float64   `json:"height" validate:"required,gt=0"`
	Unit            string    `json:"unit" validate:"required,oneof=centimeters meters feet inches"`
	Type            string    `json:"type" validate:"required,oneof=digital static LED traditional"`
	Price           float64   `json:"price" validate:"required,gt=0"`
}

type BillboardTypeRequest struct {
	Name string `json:"name" validate:"required"`
}
