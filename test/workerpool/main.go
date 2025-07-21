package main

import (
	"context"
	"ethereum-scanner/eth"
	"fmt"
	"sync"
)

// import (
// 	"ethereum-scanner/workerpool"
// 	"log"
// 	"os"
// 	"os/signal"
// 	"sync"
// 	"syscall"
// 	"time"
// )

func main() {
	// 	var wg sync.WaitGroup
	// 	ctx, cancel := context.WithCancel(context.Background())
	// 	defer cancel()

	// 	//interupt if user press ctrl+c
	// 	c := make(chan os.Signal, 1)
	// 	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	// 	go func() {
	// 		<-c
	// 		cancel()
	// 	}()

	// 	workerpool.InitWorkerPool(10, &wg)

	// 	for i := 0; i < 100; i++ {
	// 		index := i
	// 		workerpool.WorkerPool.AddJob(func() {
	// 			job(index)
	// 		})
	// 	}

	// 	<-ctx.Done()
	// 	workerpool.WorkerPool.Close()
	// 	log.Println("WorkerPool 關閉")

	// }

	//	func job(i int) {
	//		log.Printf("job %d 開始", i)
	//		time.Sleep(30 * time.Second)
	//		log.Printf("job %d 完成", i)

	eth.InitEthClient()
	defer eth.CloseEthClient()
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			blockNumber, err := eth.Client.BlockNumber(context.Background())
			if err != nil {
				fmt.Printf("failed %d: %v\n", i, err)
			} else {
				fmt.Printf("success %d: block number %d\n", i, blockNumber)
			}
		}(i)
	}

	wg.Wait()

}
