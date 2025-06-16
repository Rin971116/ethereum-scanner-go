package background

import (
	"context"
	"log"
	"math/big"
	"sync"
	"time"

	"ethereum-scanner/database"
	"ethereum-scanner/domain/entity"
	"ethereum-scanner/eth"

	"github.com/ethereum/go-ethereum/core/types"
)

func StartScanner(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done() // 確保在函數結束時減少計數器

	// 創建一個 channel 來傳遞交易和區塊號碼
	type txWithBlock struct {
		tx          *types.Transaction
		blockNumber uint64
	}
	txChan := make(chan txWithBlock, 2000)

	// 啟動一個獨立的 goroutine 來處理交易
	wg.Add(1)
	go func() {
		defer wg.Done()
		for txData := range txChan { //會阻塞，直到有交易被放入channel
			entityTx := entity.ConvertToEntityTransaction(txData.tx, txData.blockNumber)
			if entityTx == nil {
				continue
			}
			err := database.DB.Create(entityTx).Error
			if err != nil {
				log.Printf("寫入交易錯誤: %v", err)
				continue
			}
			log.Printf("成功寫入交易 %s，blockNumber: %d", txData.tx.Hash().Hex(), txData.blockNumber)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			log.Println("收到關閉信號，準備結束background scanner...")
			close(txChan) // 關閉 存放transaction的channel
			return
		default:
			// 取得最新block number
			blockNumber, err := eth.Client.BlockNumber(context.Background())
			if err != nil {
				log.Printf("取得區塊號碼錯誤: %v，等待 5 秒後重試", err)
				time.Sleep(5 * time.Second)
				continue
			} else {
				log.Printf("success 取得最新區塊號碼: %d", blockNumber)
			}

			// 檢查區塊是否已存在database
			var result entity.Block
			err = database.DB.Table("blocks").Where("number = ?", blockNumber).First(&result).Error
			if err == nil {
				log.Printf("區塊 %d 已存在，等待新區塊", blockNumber)
				time.Sleep(5 * time.Second)
				continue
			}

			// 取得最新block
			block, err := eth.Client.BlockByNumber(context.Background(), big.NewInt(int64(blockNumber)))
			if err != nil {
				log.Printf("取得區塊錯誤: %v，等待 5 秒後重試", err)
				time.Sleep(5 * time.Second)
				continue
			}

			// 將block存入database
			entityBlock := entity.ConvertToEntityBlock(block)
			err = database.DB.Create(entityBlock).Error
			if err != nil {
				log.Printf("寫入區塊錯誤: %v", err)
				continue
			}
			log.Printf("成功寫入區塊 %d 到資料庫", entityBlock.Number)

			// 將交易放入 存放transaction的channel，如果channel滿了，會阻塞，直到有空間可以放入
			for _, tx := range block.Transactions() {
				txChan <- txWithBlock{tx: tx, blockNumber: blockNumber}
			}

			// 成功處理完一個區塊後，等待一下再處理下一個
			time.Sleep(1 * time.Second)
		}
	}
}
