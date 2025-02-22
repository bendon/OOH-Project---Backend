package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BillboardSequenceModel struct {
	ID             uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	OrganizationId uuid.UUID      `gorm:"type:char(36);column:organization_id;null" json:"organizationId"`
	BoardNumber    int64          `gorm:"type:bigint;column:board_number;null" json:"boardNumber"`
	CreatedAt      int64          `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt      int64          `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (BillboardSequenceModel) TableName() string {
	return "bill_board_sequences"
}

func (b *BillboardSequenceModel) BeforeCreate(tx *gorm.DB) (err error) {
	// var maxBoardNumber int64
	// err = tx.Model(&BillboardSequenceModel{}).
	// 	Select("COALESCE(MAX(board_number), 0)").
	// 	Where("organization_id = ?", b.OrganizationId).
	// 	Scan(&maxBoardNumber).Error
	// if err != nil {
	// 	return err
	// }
	// b.BoardNumber = maxBoardNumber + 1
	b.ID = uuid.New()
	b.CreatedAt = time.Now().Unix()
	b.UpdatedAt = time.Now().Unix()
	return
}

func (b *BillboardSequenceModel) BeforeUpdate(tx *gorm.DB) (err error) {

	b.UpdatedAt = time.Now().Unix()
	return

}
