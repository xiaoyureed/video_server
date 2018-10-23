package taskrunner

const (
	READY_TO_DISPATCH = "d" // 标识准备好 dispatch了
	READY_TO_EXECUTE  = "e" // 标识准备好 exec 了
	CLOSE             = "c" // 出错, 关闭channel

	// VIDEO_PATH = ""
)

type controlChan chan string

type dataChan chan interface{}

// dispatcher, executor 都是基于这个类型
type fn func(dc dataChan) error
