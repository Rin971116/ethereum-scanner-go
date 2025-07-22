package database

import (
	"ethereum-scanner/domain/entity"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitPostgres() {
	// 可從環境變數讀取，這裡用硬編寫死的連線字串做範例
	dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable"

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("無法連接資料庫:", err)
	}

	// 自動更新資料表結構，讓資料庫中的表跟entity中定義的結構一致，如果沒有該表的話，還會自動建立該表
	err = DB.AutoMigrate(&entity.Block{}, &entity.Transaction{}, &entity.FailedTransaction{})
	if err != nil {
		log.Fatal("自動遷移資料庫結構失敗:", err)
	}

	fmt.Println("成功連接 PostgreSQL")
}

func ClosePostgres() {
	sqlDB, err := DB.DB()
	if err != nil {
		log.Printf("無法獲取資料庫實例: %v", err)
		return
	}
	sqlDB.Close()
	log.Println("資料庫連線已關閉")
}
