// 变体 B：每任务一 goroutine，但用“带缓冲 channel 当信号量”限制同时运行的个数。
//
// 特点：
//   - 既保留“每任务一 goroutine”的简单分发，又把并发数限制在 maxConcurrent（5）以内。
//   - 适合“任务数很多、但只想同时跑 N 个”的场景（如限流调下游 API）。
//   - sem <- struct{}{} 获取配额（满了就阻塞），defer <-sem 释放配额。
//
// 运行：go run -race ./week01-go-concurrency/experiments/concurrent-print/semaphore/
package main

import (
	"fmt"
	"sync"
)

const (
	numTasks      = 100
	maxConcurrent = 5
)

func main() {
	var wg sync.WaitGroup
	sem := make(chan struct{}, maxConcurrent)

	for i := 1; i <= numTasks; i++ {
		wg.Add(1)
		sem <- struct{}{} // 获取一个并发配额；满了则阻塞，直到有 goroutine 释放
		go func(id int) {
			defer wg.Done()
			defer func() { <-sem }() // 任务结束，释放配额
			fmt.Printf("任务 #%d 完成\n", id)
		}(i)
	}
	wg.Wait()
	fmt.Println("全部任务完成")
}
