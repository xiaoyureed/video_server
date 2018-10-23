package taskrunner

// Runner run video deletion task
type Runner struct {
	Controller controlChan
	Error      controlChan
	Data       dataChan
	dataSize   int
	longLived  bool // 是否长期存活， 决定是否回收
	Dispatcher fn
	Executor   fn
}

// NewRunner - constructor
func NewRunner(size int, longlived bool, d fn, e fn) *Runner {
	return &Runner{
		Controller: make(chan string, 1),
		Error:      make(chan string, 1),
		Data:       make(chan interface{}, size),
		longLived:  longlived,
		dataSize:   size,
		Dispatcher: d,
		Executor:   e,
	}
}

func (r *Runner) startDispatch() {
	defer func() {
		if !r.longLived {
			close(r.Controller)
			close(r.Data)
			close(r.Error)
		}
	}()

	for { // 死循环, 表示 常驻任务
		select { // 监控多个 channel, 有数据流入, 则取出来
		case c := <-r.Controller:
			if c == READY_TO_DISPATCH {
				err := r.Dispatcher(r.Data)
				if err != nil {
					r.Error <- CLOSE
				} else {
					r.Controller <- READY_TO_EXECUTE
				}
			}

			if c == READY_TO_EXECUTE {
				err := r.Executor(r.Data)
				if err != nil {
					r.Error <- CLOSE
				} else {
					r.Controller <- READY_TO_DISPATCH
				}
			}
		case e := <-r.Error:
			if e == CLOSE {
				return
			}
		default:

		}
	}
}

// Start 启动 runner
func (r *Runner) Start() {
	r.Controller <- READY_TO_DISPATCH
	r.startDispatch()
}
