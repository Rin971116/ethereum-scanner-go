package entity

import (
	"github.com/ethereum/go-ethereum/core/types"
	"gorm.io/gorm"
)

type Block struct {
	gorm.Model // 自動帶入 ID、CreatedAt、UpdatedAt、DeletedAt

	Number     uint64 `gorm:"uniqueIndex;not null"` // 會創建 number 列，設置為唯一索引
	Hash       string `gorm:"size:66;not null"`     // 會創建 hash 列，長度限制為 66
	ParentHash string `gorm:"size:66"`              // 會創建 parent_hash 列
	Timestamp  uint64 `gorm:"not null"`             // 會創建 timestamp 列
	Miner      string `gorm:"size:42"`              // 會創建 miner 列
	GasLimit   uint64 // 會創建 gas_limit 列
	GasUsed    uint64 // 會創建 gas_used 列
	BaseFee    string `gorm:"type:text"` // 會創建 base_fee 列，類型為 text
	TxCount    int    `gorm:"not null"`  // 會創建 tx_count 列
}

// 將ethereum返回的types.Block轉換為我們自定義的Block，只保留我們需要的字段
func ConvertToEntityBlock(block *types.Block) *Block {
	return &Block{
		// 只轉換我們需要的字段
		Number:     block.Number().Uint64(),
		Hash:       block.Hash().Hex(),
		ParentHash: block.ParentHash().Hex(),
		Timestamp:  block.Time(),
		Miner:      block.Coinbase().Hex(),
		GasLimit:   block.GasLimit(),
		GasUsed:    block.GasUsed(),
		BaseFee:    block.BaseFee().String(),
		TxCount:    len(block.Transactions()),
		// 其他字段我們不需要，直接忽略
	}
}
