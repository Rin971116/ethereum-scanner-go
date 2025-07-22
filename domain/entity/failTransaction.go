package entity

import (
	"gorm.io/gorm"
)

type FailedTransaction struct {
	gorm.Model
	TxHash      string `gorm:"uniqueIndex;type:varchar(66)"`
	BlockNumber uint64 `gorm:"index"`
	Reason      string `gorm:"type:varchar(255)"`
}
