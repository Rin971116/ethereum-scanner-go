package handler

import (
	"ethereum-scanner/database"
	"ethereum-scanner/domain/entity"
	"log"
	"net/http"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/gin-gonic/gin"
)

func GetTransactionHandler(c *gin.Context) {
	// 1️⃣ 從 URL 的路徑參數中取得 "blockNumber" 的值
	blockNumberStr := c.Param("blockNumber")

	// 2️⃣ 將取得的 number 由 string 轉成 uint64
	blockNumber, err := strconv.ParseUint(blockNumberStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid block number"})
		return
	}

	// 3️⃣ 從資料庫查詢該區塊的交易
	var transactions []entity.Transaction
	result := database.DB.Where("block_number = ?", blockNumber).Find(&transactions)
	if result.Error != nil {
		log.Println("查詢交易錯誤:", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transactions"})
		return
	}

	// 4️⃣ 回傳交易資料
	c.IndentedJSON(http.StatusOK, transactions)
}

type TxWithFrom struct {
	Transaction *types.Transaction `json:"transaction"`
	From        common.Address     `json:"from"`
}
