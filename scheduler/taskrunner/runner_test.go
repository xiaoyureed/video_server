package taskrunner

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func TestRunner(t *testing.T) {
	// t.SkipNow()
	dispatch := func(dc dataChan) error {
		for i := 0; i < 10; i++ {
			dc <- i
			log.Printf("Dispatcher sent: %v", i)
		}
		return nil
	}

	exec := func(dc dataChan) error {
	forloop:
		for {
			select {
			case data := <-dc:
				log.Printf("Executor received: %v", data)
			default:
				break forloop
			}
		}
		return nil
		// return errors.New("Executor")
	}

	runner := NewRunner(30, false, dispatch, exec)
	go runner.Start()
	time.Sleep(3 * time.Second)
}

func TestDemo(t *testing.T) {
	t.SkipNow()
	go func() {
		for {
			fmt.Println("aaaaaaaaaaaaaaaa")
		}
	}()

	time.Sleep(1 * time.Second)
}
