package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"bbscout/types"
)

type BillboardModel struct {
	ID              uuid.UUID               `gorm:"type:char(36);primaryKey" json:"id"`
	OrganizationId  uuid.UUID               `gorm:"type:char(36);column:organization_id;null" json:"organizationId"`
	CreatedById     uuid.UUID               `gorm:"type:char(36);column:created_by_id; null" json:"createdById"`
	BoardCode       string                  `gorm:"type:varchar(1000);column:board_code; null" json:"boardCode"`
	Owner           *string                 `gorm:"type:varchar(255);column:owner; null" json:"Owner"`
	OwnerContacts   *types.Int64ArrayJSONB  `gorm:"type:json;column:owner_contact;null" json:"ownerContacts"`
	OwnerEmails     *types.StringArrayJSONB `gorm:"type:json;column:owner_email;null" json:"ownerEmails"`
	Description     *string                 `gorm:"type:text;column:description;" json:"description"`
	Location        *string                 `gorm:"type:varchar(255);column:location;null" json:"location"`
	City            *string                 `gorm:"type:varchar(255);column:city;null" json:"city"`
	Latitude        float64                 `gorm:"type:double;column:latitude;not null" json:"latitude"`
	Longitude       float64                 `gorm:"type:double;column:longitude;not null" json:"longitude"`
	Accuracy        *float64                `gorm:"type:double;column:accuracy;null; default:0" json:"accuracy"`
	Width           *float64                `gorm:"type:double;column:width;not null;default:0" json:"width"`
	Height          *float64                `gorm:"type:double;column:height;not null;default:0" json:"height"`
	Unit            string                  `gorm:"type:enum('centimeters','meters','feet','inches');column:unit;not null;default:meters" json:"unit"`
	Type            string                  `gorm:"type:varchar(255);column:type;not null" json:"type"` // "Static Billboard", "Digital Billboard", "Banner Ads", "Wallscapes", "Mobile Billboards","Lamp Posts","Interactive Billboards"
	ObjectType      *string                 `gorm:"type:varchar(255);column:object_type;null" json:"objectType"`
	ParentBoardCode *string                 `gorm:"type:varchar(1000);column:parent_board_code;null" json:"parentBoardCode"`
	Price           *float64                `gorm:"type:decimal(10,2);column:price; null; default:0" json:"price"`
	ImageId         *uuid.UUID              `gorm:"type:char(36);column:image_id; null" json:"imageId"`
	CloseUpImageId  *uuid.UUID              `gorm:"type:char(36);column:closeup_image_id; null" json:"closeUpImageId"`
	Image           *FileModel              `gorm:"foreignKey:ImageId; references:ID" json:"image,omitempty"`
	CloseUpImage    *FileModel              `gorm:"foreignKey:CloseUpImageId; references:ID" json:"closeUpImage,omitempty"`
	Active          bool                    `gorm:"type:boolean;column:active;default:true;not null" json:"active"`
	Occupied        bool                    `gorm:"type:boolean;column:occupied;default:false;not null" json:"occupied"`
	Structure       *string                 `gorm:"type:varchar(255);column:structure; null" json:"structure"`
	Material        *string                 `gorm:"type:varchar(255);column:material; null" json:"material"`
	Angle           *string                 `gorm:"type:varchar(255);column:angel; null" json:"angle"`
	Visibility      *string                 `gorm:"type:varchar(255);column:visibility; null" json:"visibility"`
	Illumination    *string                 `gorm:"type:varchar(255);column:illumination;null" json:"illumination"`
	CreatedAt       int64                   `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt       int64                   `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt       gorm.DeletedAt          `gorm:"index" json:"-"`
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
