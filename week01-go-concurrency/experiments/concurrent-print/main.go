// 实验：100 个任务由多个 goroutine 并发执行，每个任务打印自己的编号。
//
// 演示要点：
//  1. 用固定数量的 worker goroutine 从一个 channel 消费 100 个任务（多个 goroutine 共享工作）。
//  2. 用 sync.WaitGroup 在所有任务完成后才退出 main，绝不依赖 time.Sleep 假设完成。
//  3. 多个 goroutine 同时打印时，输出顺序不确定（详见下方“为什么顺序不确定”说明）。
//
// 运行：
//   go run -race ./week01-go-concurrency/experiments/concurrent-print/
package main

import (
	"fmt"
	"sync"
)

const (
	numTasks   = 100 // 任务总数
	numWorkers = 5   // 并发 worker 数量（远小于任务数，体现“多 goroutine 消费”）
)

func main() {
	taskCh := make(chan int)
	var wg sync.WaitGroup

	// 启动 numWorkers 个 worker goroutine，每个持续从 taskCh 取任务执行。
	for w := 0; w < numWorkers; w++ {
		go func(workerID int) {
			for taskID := range taskCh {
				// 每个任务打印自己的编号，并标注是哪个 worker 执行的。
				fmt.Printf("worker %d 执行任务 #%d\n", workerID, taskID)
				wg.Done() // 当前任务完成，计数器 -1。
			}
		}(w)
	}

	// 在主 goroutine 中提交任务：先 Add(1) 再发送，保证 Add 一定发生在对应 Done 之前。
	for i := 1; i <= numTasks; i++ {
		wg.Add(1)
		taskCh <- i
	}

	close(taskCh) // 通知所有 worker：没有更多任务了。
	wg.Wait()     // 阻塞，直到 100 个 Done() 全部执行完毕。

	fmt.Println("全部任务完成")
}
