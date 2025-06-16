package main

import (
	"ethereum-scanner/background"
	"ethereum-scanner/controllers"
	"ethereum-scanner/database"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"context"
	"ethereum-scanner/eth"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	var wg sync.WaitGroup

	// 建立 context，用於優雅關閉伺服器
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel() //關閉 負責監聽中斷操作的 gorutine

	// 初始化以太坊客戶端
	eth.InitEthClient()

	// 初始化資料庫連線
	database.InitPostgres()

	// 起一個goroutine 跑背景掃鏈
	wg.Add(1)
	go background.StartScanner(ctx, &wg)

	// 起一個goroutine 跑 gin Web Server（避免阻塞主線程）
	wg.Add(1)
	go StartGinServer(ctx, &wg)

	// 阻塞在這邊直到程式接收到中斷訊號（如 Ctrl+C）
	<-ctx.Done()
	log.Println("main收到中斷訊號，等待其他 goroutine 結束...")

	// 等待所有 goroutine 正常結束
	wg.Wait()
	log.Println("所有 goroutine 結束，程式關閉")

	// 關閉 Ethereum client（釋放 RPC 資源）
	if eth.Client != nil {
		eth.Client.Close()
		log.Println("Ethereum client closed.")
	}

	// 關閉資料庫連線
	sqlDB, err := database.DB.DB()
	if err != nil {
		log.Printf("獲取資料庫實例錯誤: %v", err)
	} else {
		err := sqlDB.Close()
		if err != nil {
			log.Printf("關閉資料庫連線錯誤: %v", err)
		}
	}
}

func StartGinServer(ctx context.Context, wg *sync.WaitGroup) {

	// goroutine 結束時 wg-1
	defer wg.Done()

	// 啟動 Gin 伺服器
	router := gin.Default()

	// 註冊 API 路由
	router.GET("/hello", controllers.HelloWorldHandler)
	router.GET("/block/:blockNumber", controllers.GetBlockHandler)
	router.GET("/transaction/:hash", controllers.GetTransactionHandler)
	// router.GET("/eventlog", controllers.GetEventLogHandler)
	// router.POST("/migrate", controllers.MigrateHandler)

	// 包成 http.Server，能更靈活控制
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// 啟動 HTTP server（非阻塞）
	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Gin server 啟動失敗: %v\n", err)
		}
	}()

	log.Println("Gin Server 已啟動在 http://localhost:8080")

	// 等待 ctx 結束信號（例如 ctrl+c）
	<-ctx.Done()
	log.Println("收到中斷訊號，正在關閉 Gin Server...")

	// 嘗試在 5 秒內優雅關閉 HTTP server
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := srv.Shutdown(shutdownCtx)
	if err != nil {
		log.Fatalf("Gin Server 無法優雅關閉: %v", err)
	}

	log.Println("Gin Server 已成功關閉")
}
