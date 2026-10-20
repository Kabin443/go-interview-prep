// Package workerpool 实现了一个固定数量的 goroutine 工作池，
// 用于并发消费任务（Task）。这是 Go 并发面试中的经典考点，
// 涉及 goroutine、channel、sync.WaitGroup 以及 context 取消控制。
package workerpool

import (
	"context"
	"sync"
)

// Task 表示一个可被 worker 执行的工作单元。
// 它接收 context 以便支持取消，并返回执行结果（error）。
type Task func(ctx context.Context) error

// WorkerPool 管理一组长期运行的 goroutine，
// 它们从同一个任务 channel 中竞争消费任务。
type WorkerPool struct {
	tasks   chan Task
	workers int

	wg   sync.WaitGroup
	mu   sync.Mutex
	errs []error
}

// New 创建一个指定 worker 数量的 WorkerPool。
// 如果 workers <= 0，则退化为单个 worker，保证至少有一个消费者。
func New(workers int) *WorkerPool {
	if workers <= 0 {
		workers = 1
	}
	return &WorkerPool{
		tasks:   make(chan Task),
		workers: workers,
	}
}

// worker 是长期运行的 goroutine：不断从任务 channel 取任务执行，
// 直到任务 channel 被关闭，或传入的 context 被取消。
func (p *WorkerPool) worker(ctx context.Context) {
	defer p.wg.Done()
	for {
		select {
		case <-ctx.Done():
			// 外部取消信号，立即退出，不再处理新任务。
			return
		case task, ok := <-p.tasks:
			if !ok {
				// channel 已关闭，所有任务已消费完毕。
				return
			}
			if err := task(ctx); err != nil {
				p.mu.Lock()
				p.errs = append(p.errs, err)
				p.mu.Unlock()
			}
		}
	}
}

// Start 启动所有 worker goroutine。必须在 Submit 之前调用。
func (p *WorkerPool) Start(ctx context.Context) {
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go p.worker(ctx)
	}
}

// Submit 提交一个任务。如果所有 worker 都忙碌，调用会阻塞直到有空位。
func (p *WorkerPool) Submit(task Task) {
	p.tasks <- task
}

// Stop 关闭任务 channel 并等待所有 worker 退出。
// 调用 Stop 后应不再 Submit，否则会 panic（向已关闭的 channel 发送）。
func (p *WorkerPool) Stop() {
	close(p.tasks)
	p.wg.Wait()
}

// Errors 返回已收集的任务执行错误副本。
func (p *WorkerPool) Errors() []error {
	p.mu.Lock()
	defer p.mu.Unlock()
	out := make([]error, len(p.errs))
	copy(out, p.errs)
	return out
}
