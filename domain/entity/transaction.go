package entity

import (
	"context"
	"ethereum-scanner/eth"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"gorm.io/gorm"
)

// Transaction 交易實體
type Transaction struct {
	gorm.Model
	Hash        string `gorm:"uniqueIndex;type:varchar(66)"` // 交易primary哈希值
	BlockNumber uint64 `gorm:"index"`                        // 區塊號碼
	From        string `gorm:"type:varchar(42)"`             // 發送方地址
	To          string `gorm:"type:varchar(42)"`             // 接收方地址
	Value       string `gorm:"type:varchar(32)"`             // 交易金額
	GasPrice    string `gorm:"type:varchar(32)"`             // Gas 價格
	GasLimit    uint64 // Gas 限制
	GasUsed     uint64 // 使用的 Gas
	Nonce       uint64 // 交易序號
	Data        []byte `gorm:"type:bytea"` // 交易數據
	Status      bool   // 交易狀態（成功/失敗）
}

// ConvertToEntityTransaction 將以太坊交易轉換為實體
func ConvertToEntityTransaction(tx *types.Transaction, blockNumber uint64) *Transaction {
	// 獲取交易收據以確認狀態(status)
	receipt, err := eth.Client.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		log.Printf("獲取交易收據錯誤: %v", err)
		return nil
	}

	// 獲取發送方地址(from)
	signer := types.LatestSignerForChainID(big.NewInt(1))
	from, err := signer.Sender(tx)
	if err != nil {
		log.Printf("獲取發送方地址錯誤: %v", err)
		return nil
	}

	// 獲取接收方地址(to)，如果該transaction沒有接收方地址(像是合約創建交易)，則to為空字符串
	var to string
	if tx.To() != nil {
		to = tx.To().Hex()
	}

	return &Transaction{
		Hash:        tx.Hash().Hex(),
		BlockNumber: blockNumber,
		From:        from.Hex(),
		To:          to,
		Value:       tx.Value().String(),
		GasPrice:    tx.GasPrice().String(),
		GasLimit:    tx.Gas(),
		GasUsed:     receipt.GasUsed,
		Nonce:       tx.Nonce(),
		Data:        tx.Data(),
		Status:      receipt.Status == 1,
	}
}
