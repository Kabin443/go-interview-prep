// 变体 A：每个任务单独起一个 goroutine，全部交给 WaitGroup 等待。
//
// 特点：
//   - 写法最简单直观，goroutine 数量 = 任务数（这里是 100 个）。
//   - 任务量很大时（如百万级）会创建过多 goroutine，需改用 worker 池或信号量限流。
//   - 注意：循环变量 id 按值传入闭包，避免“闭包捕获同一变量”的经典陷阱
//     （Go 1.22+ 已默认每轮新变量，但按值传入仍是稳妥写法）。
//
// 运行：go run -race ./week01-go-concurrency/experiments/concurrent-print/one-goroutine-per-task/
package main

import (
	"fmt"
	"sync"
)

const numTasks = 100

func main() {
	var wg sync.WaitGroup
	for i := 1; i <= numTasks; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			fmt.Printf("任务 #%d 完成\n", id)
		}(i)
	}
	wg.Wait()
	fmt.Println("全部任务完成")
}
