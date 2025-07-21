package workerpool

import (
	"log"
	"sync"
)

type workerPool struct {
	WorkerCount int
	JobQueue    chan func()
}

// 目的是透過這個WorkerPool來一次管理多個goroutine
var WorkerPool *workerPool

func InitWorkerPool(workerCount int, wg *sync.WaitGroup) {
	WorkerPool = &workerPool{
		WorkerCount: workerCount,
		JobQueue:    make(chan func(), workerCount*1000),
	}

	// 起 workerCount 個 goroutine
	for i := 0; i < workerCount; i++ {
		workerId := i

		wg.Add(1)
		go func(workerId int) {
			defer wg.Done()
			for job := range WorkerPool.JobQueue { // for range 是當job channel 被關閉時，並且裡面沒有任何元素，會自動退出for loop
				log.Printf("Worker %d 拿了一個 job", workerId)
				job()
				log.Printf("Worker %d 完成了一個 job", workerId)
			}
		}(workerId)

		log.Printf("已建立%d號 worker", workerId)
	}
}

func (wp *workerPool) AddJob(job func()) {
	wp.JobQueue <- job
	log.Println("WorkerPool 收到 job")
}

func (wp *workerPool) Close() {
	close(wp.JobQueue)
	log.Println("WorkerPool 關閉 -> job channel 已關閉，等待所有 worker 結束...")
}
