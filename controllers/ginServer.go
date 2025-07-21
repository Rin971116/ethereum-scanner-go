package controllers

import (
	"context"
	"ethereum-scanner/controllers/handler"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

func StartGinServer(ctx context.Context, wg *sync.WaitGroup) {

	// goroutine 結束時 wg-1
	defer wg.Done()

	// 啟動 Gin 伺服器
	router := gin.Default()

	// 註冊 API 路由
	router.GET("/block/:blockNumber", handler.GetBlockHandler)
	router.GET("/transaction/:hash", handler.GetTransactionHandler)
	router.GET("/eventlog", handler.GetEventLogHandler)
	router.POST("/migrate", handler.MigrateHandler)

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
