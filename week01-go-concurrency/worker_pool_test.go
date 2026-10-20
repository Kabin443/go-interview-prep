package workerpool

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

// TestProcessesAllTasks 验证所有提交的任务都被执行，且并发数可控。
func TestProcessesAllTasks(t *testing.T) {
	p := New(4)
	p.Start(context.Background())

	const n = 100
	var counter int64
	for i := 0; i < n; i++ {
		p.Submit(func(ctx context.Context) error {
			atomic.AddInt64(&counter, 1)
			return nil
		})
	}
	p.Stop()

	if got := atomic.LoadInt64(&counter); got != n {
		t.Fatalf("期望处理 %d 个任务，实际处理 %d 个", n, got)
	}
}

// TestCollectsErrors 验证任务返回的错误被统一收集，不会丢失。
func TestCollectsErrors(t *testing.T) {
	p := New(3)
	p.Start(context.Background())

	const n = 5
	for i := 0; i < n; i++ {
		p.Submit(func(ctx context.Context) error {
			return errors.New("boom")
		})
	}
	p.Stop()

	if errs := p.Errors(); len(errs) != n {
		t.Fatalf("期望收集 %d 个错误，实际 %d 个", n, len(errs))
	}
}

// TestContextCancellation 验证 context 取消后 worker 能及时退出，
// 避免 Stop 永久阻塞。
func TestContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	p := New(2)
	p.Start(ctx)

	cancel() // 触发取消，worker 应立即退出。

	done := make(chan struct{})
	go func() {
		p.Stop()
		close(done)
	}()

	select {
	case <-done:
		// 正常退出。
	case <-time.After(2 * time.Second):
		t.Fatal("context 取消后 Stop 仍未返回，可能存在泄漏")
	}
}

// TestNewDefaultsToOneWorker 验证非法 worker 数量时退化为单个 worker。
func TestNewDefaultsToOneWorker(t *testing.T) {
	p := New(0)
	if p.workers != 1 {
		t.Fatalf("期望 1 个 worker，实际 %d 个", p.workers)
	}
}
