package taskrunner

import (
	"errors"
	"log"
	"os"
	"sync"
	"video_server/common"
	"video_server/scheduler/db"
)

// deleteVideo delete real video file
func deleteVideo(vid string) error {
	err := os.Remove(common.VIDEO_DIR + vid)

	if err != nil && !os.IsNotExist(err) {
		log.Printf("Deleting video error: %v", err)
		return err
	}

	return nil
}

// VideoClearDispatcher dispatch video ids which would be deleted to dataChan
func VideoClearDispatcher(dc dataChan) error {
	ids, err := db.ReadVideoDeletionRecord(3) // read 3 deletion record
	if err != nil {
		log.Printf("Video clear dispatcher error: %v", err)
		return err
	}

	if len(ids) == 0 {
		return errors.New("All tasks finished")
	}

	for _, id := range ids {
		dc <- id
	}

	return nil
}

//VideoClearExecutor receive video ids from dataChan and delete video
func VideoClearExecutor(dc dataChan) error {
	errMap := &sync.Map{}
	var err error

forloop:
	for {
		select {
		case vid := <-dc:
			go func(id interface{}) {
				if err := deleteVideo(id.(string)); err != nil {
					errMap.Store(id, err)
					return
				}
				if err := db.DelVideoDeletionRecord(id.(string)); err != nil {
					errMap.Store(id, err)
					return
				}
			}(vid)
		default:
			break forloop
		}
	}

	errMap.Range(func(k, v interface{}) bool {
		err = v.(error)
		if err != nil {
			return false // 只要有 error , 就停止迭代
		}
		return true
	})

	return err
}
