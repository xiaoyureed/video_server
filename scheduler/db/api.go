package db

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
)

//AddVideoDeletionRecord add video id which to be deleted
func AddVideoDeletionRecord(vid string) error {
	stmtIns, err := conn.Prepare("INSERT INTO video_to_del (video_id) VALUES(?)")
	if err != nil {
		return err
	}

	_, err = stmtIns.Exec(vid)
	if err != nil {
		log.Printf("AddVideoDeletionRecord error: %v", err)
		return err
	}

	defer stmtIns.Close()
	return nil
}

//ReadVideoDeletionRecord 批量获取 待删除 video id
func ReadVideoDeletionRecord(count int) ([]string, error) {
	stmtOut, err := conn.Prepare("SELECT video_id FROM video_to_del LIMIT ?")

	var ids []string

	if err != nil {
		return ids, err
	}

	rows, err := stmtOut.Query(count)
	if err != nil {
		log.Printf("Query VideoDeletionRecord error: %v", err)
		return ids, err
	}

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return ids, err
		}

		ids = append(ids, id)
	}

	defer stmtOut.Close()
	return ids, nil
}

//DelVideoDeletionRecord
func DelVideoDeletionRecord(vid string) error {
	stmtDel, err := conn.Prepare("DELETE FROM video_to_del WHERE video_id=?")
	if err != nil {
		return err
	}

	_, err = stmtDel.Exec(vid)
	if err != nil {
		log.Printf("Deleting VideoDeletionRecord error: %v", err)
		return err
	}

	defer stmtDel.Close()
	return nil
}
