package controllers

import (
	"ethereum-scanner/database"
	"ethereum-scanner/domain/entity"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetBlockHandler(c *gin.Context) {
	// 1️⃣ 從 URL 的路徑參數中取得 "number" 的值
	blockNumberStr := c.Param("blockNumber")

	// 2️⃣ 將取得的 number 由 string 轉成 uint64
	blockNumber, err := strconv.ParseUint(blockNumberStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid block number"})
		return
	}

	// 3️⃣ 從資料庫查詢該區塊
	var block entity.Block
	result := database.DB.Select("blocks").Where("number = ?", blockNumber).First(&block)
	if result.Error != nil {
		log.Println("查詢區塊錯誤:", result.Error)
		c.JSON(http.StatusNotFound, gin.H{"error": "Block not found"})
		return
	}

	// 4️⃣ 回傳區塊資料
	c.IndentedJSON(http.StatusOK, block)
}

type BlockWithTransactions struct {
	Block        *entity.Block        `json:"block"`
	Transactions []entity.Transaction `json:"transactions"`
}
