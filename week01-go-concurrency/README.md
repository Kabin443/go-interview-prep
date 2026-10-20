# Week 01 · Go 并发（Concurrency）

本周聚焦 Go 并发模型的核心概念与典型模式，所有代码与实验均放在本目录。

## 目录结构

```
week01-go-concurrency/
├── README.md              # 本周说明（本文件）
├── worker_pool.go         # 工作池实现
├── worker_pool_test.go    # 工作池单元测试
└── experiments/           # 各类并发小实验（见下）
```

## 学习目标

- 理解 goroutine 的调度模型与开销（M:N 调度、GPM 模型）。
- 掌握 channel 的几种用法：无缓冲/有缓冲、关闭语义、select 多路复用。
- 掌握 `sync.WaitGroup`、`sync.Mutex` 在并发协作与共享状态保护中的作用。
- 理解 `context` 在取消传播、超时控制中的角色。
- 能独立实现一个健壮的 worker pool，并解释其并发安全点。

## worker_pool.go 要点

- 固定数量的 worker goroutine，从同一个 `tasks` channel 竞争消费任务。
- `Start` 启动 worker；`Submit` 提交任务（worker 全忙时阻塞）；`Stop` 关闭 channel 并等待退出。
- 通过 `context` 支持外部取消；多 goroutine 写错误切片时用 `sync.Mutex` 保护。
- 关键并发安全点：`errs` 切片的追加受锁保护；`Stop` 关闭 channel 后不再 `Submit`。

## 运行测试

```bash
# 在项目根目录（已 go mod init go-interview-prep）
go test ./week01-go-concurrency/...

# 带竞态检测，强烈建议每次都开 -race
go test -race ./week01-go-concurrency/...
```

## experiments/ 说明

该目录用于存放本周的各类小实验，例如：
- `buffered-vs-unbuffered`：缓冲/无缓冲 channel 行为对比。
- `select-timeout`：select + time.After 实现超时。
- `context-propagation`：取消信号的跨 goroutine 传播。
- `deadlock-demo`：故意制造并分析死锁。

每个实验建议单独建子目录或 `.go` 文件，并附带简短注释说明观察结论。
（当前为占位目录，提交一个 `.gitkeep` 以便纳入版本控制。）
