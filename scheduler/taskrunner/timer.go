package taskrunner

import "time"

/*
逻辑: 创建&启动定时器 Worker <- 创建&启动任务运行器 runner
*/

//Worker 定时器
type Worker struct {
	ticker *time.Ticker
	runner *Runner
}

// NewWorker constructor
func NewWorker(interval time.Duration, r *Runner) *Worker {
	return &Worker{
		ticker: time.NewTicker(interval * time.Second), // 为 ticker 设置 间隔
		runner: r,
	}
}

//startWorker 启动定时器
func (w *Worker) startWorker() {
	for {
		select {
		// 每隔段时间会有数据流动, 这就形成了定时功能
		case <-w.ticker.C: // 不要使用 for range 取ticker, 是同步的
			go w.runner.Start()
		}
	}
}

// Start 定时器启动, 整个 task runner 启动
func Start() {
	// Start video file cleaning
	r := NewRunner(3, true, VideoClearDispatcher, VideoClearExecutor)
	w := NewWorker(3, r)
	go w.startWorker()
}
