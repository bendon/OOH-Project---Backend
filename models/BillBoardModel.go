package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BillboardModel struct {
	ID              uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	OrganizationId  uuid.UUID      `gorm:"type:char(36);column:organization_id;null" json:"organizationId"`
	CreatedById     uuid.UUID      `gorm:"type:char(36);column:created_by_id; null" json:"createdById"`
	BoardCode       string         `gorm:"type:varchar(1000);column:board_code; null" json:"boardCode"`
	Title           string         `gorm:"type:varchar(255);column:title;not null" json:"title"`
	Description     string         `gorm:"type:text;column:description;" json:"description"`
	Location        string         `gorm:"type:varchar(255);column:location;not null" json:"location"`
	Latitude        float64        `gorm:"type:double;column:latitude;not null" json:"latitude"`
	Longitude       float64        `gorm:"type:double;column:longitude;not null" json:"longitude"`
	Accuracy        float64        `gorm:"type:double;column:accuracy;not null" json:"accuracy"`
	Width           float64        `gorm:"type:double;column:width;not null" json:"width"`
	Height          float64        `gorm:"type:double;column:height;not null" json:"height"`
	Unit            string         `gorm:"type:enum('centimeters','meters','feet','inches');column:unit;not null" json:"unit"`
	Type            string         `gorm:"type:enum('Static Billboard','Digital Billboard','Banner Ads','Wallscapes','Mobile Billboards','Lamp Posts','Interactive Billboards');column:type;not null" json:"type"`
	ParentBoardCode *string        `gorm:"type:varchar(1000);column:parent_board_code;null" json:"parentBoardCode"`
	Price           *float64       `gorm:"type:decimal(10,2);column:price; null; default:0" json:"price"`
	ImageId         *uuid.UUID     `gorm:"type:char(36);column:image_id; null" json:"image_id"`
	Image           *FileModel     `gorm:"foreignKey:ImageId; references:ID" json:"image,omitempty"`
	Active          bool           `gorm:"type:boolean;column:active;default:true;not null" json:"active"`
	CreatedAt       int64          `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt       int64          `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

func (BillboardModel) TableName() string {
	return "bill_boards"
}

func (b *BillboardModel) BeforeCreate(tx *gorm.DB) (err error) {
	b.ID = uuid.New()
	b.CreatedAt = time.Now().Unix()
	b.UpdatedAt = time.Now().Unix()
	return
}

func (b *BillboardModel) BeforeUpdate(tx *gorm.DB) (err error) {
	b.UpdatedAt = time.Now().Unix()
	return

}
