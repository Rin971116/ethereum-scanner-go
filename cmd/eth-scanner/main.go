package main

import (
	"context"
	"ethereum-scanner/controllers"
	"ethereum-scanner/database"
	"ethereum-scanner/eth"
	"ethereum-scanner/kafka"
	"ethereum-scanner/scanner"
	"ethereum-scanner/workerpool"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {

	var wg sync.WaitGroup

	// 建立 context，用於優雅關閉伺服器
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel() //離開前關閉 負責監聽中斷操作的 gorutine

	// 初始化以太坊客戶端
	eth.InitEthClient()
	// 離開前關閉 Ethereum client（釋放 RPC 資源）
	defer eth.CloseEthClient()

	// 初始化全域限流器，沒有用到資源，所以不用關閉
	eth.InitReceiptRateLimiter()

	// 初始化 kafka writer
	kafka.InitKafkaWriter([]string{"localhost:9092"}, "all_tx_topic")
	defer kafka.CloseKafkaWriter()

	// 初始化資料庫連線
	database.InitPostgres()
	// 離開前關閉資料庫連線
	defer database.ClosePostgres()

	// 初始化 worker pool
	workerpool.InitWorkerPool(20, &wg)

	// 起一個goroutine 跑背景掃鏈
	wg.Add(1)
	go scanner.StartScanner(ctx, &wg)

	// 起一個goroutine 跑 gin Web Server（避免阻塞主線程）-
	wg.Add(1)
	go controllers.StartGinServer(ctx, &wg)

	// 阻塞在這邊直到程式接收到中斷訊號（如 Ctrl+C）
	<-ctx.Done()
	log.Println("main收到中斷訊號，等待其他 goroutine 結束...")

	// 等待所有 goroutine 正常結束
	wg.Wait()
	log.Println("所有 goroutine 結束，程式關閉")

}
