package scanner

import (
	"context"
	"log"
	"math/big"
	"sync"
	"time"

	"ethereum-scanner/database"
	"ethereum-scanner/domain/entity"
	"ethereum-scanner/eth"
	"ethereum-scanner/kafka"
	"ethereum-scanner/workerpool"
)

func StartScanner(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done() // 確保在函數結束時減少計數器

	for {
		select {
		case <-ctx.Done():
			log.Println("scanner收到關閉信號，準備結束...")
			workerpool.WorkerPool.Close() // 關閉 存放transaction的channel
			return
		default:
			// 取得最新block number
			blockNumber, err := eth.GetNextEthClient().BlockNumber(context.Background())
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
			block, err := eth.GetNextEthClient().BlockByNumber(context.Background(), big.NewInt(int64(blockNumber)))
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
				currentTx := tx //!!每次都聲明一個新的變數(記憶體位置不同)，避免在多個goroutine中使用相同的變數(相同的記憶體位置)
				workerpool.WorkerPool.AddJob(func() {
					entityTx := entity.ConvertToEntityTransaction(currentTx, blockNumber)
					if entityTx == nil {
						log.Printf("交易 %s 轉換失敗，放棄處理", currentTx.Hash().Hex())
						database.DB.Create(&entity.FailedTransaction{
							TxHash:      currentTx.Hash().Hex(),
							BlockNumber: blockNumber,
							Reason:      "轉換失敗",
						})
						return
					}
					// 將交易放入kafka
					err = kafka.Produce(entityTx)
					if err != nil {
						log.Printf("交易寫入kafka錯誤: %v", err)
					} else {
						log.Printf("成功寫入交易 %s 到kafka", entityTx.Hash)
					}

					// 將交易寫入DB
					err = database.DB.Create(entityTx).Error
					if err != nil {
						log.Printf("交易寫入DB錯誤: %v", err)
					} else {
						log.Printf("成功寫入交易 %s 到資料庫", entityTx.Hash)
					}
				})
			}

			// 成功處理完一個區塊後，等待一下再處理下一個
			time.Sleep(1 * time.Second)
		}
	}
}
