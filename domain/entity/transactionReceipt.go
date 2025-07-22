package entity

import "gorm.io/gorm"

type TransactionReceipt struct {
	gorm.Model
	TxHash      string `gorm:"uniqueIndex;type:varchar(66)"`
	BlockNumber uint64 `gorm:"index"`
	GasUsed     uint64
	Status      bool
}
