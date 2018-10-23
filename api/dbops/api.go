package dbops

import (
	"database/sql"
	"log"
	"time"
	"video_server/api/defs"
	"video_server/api/utils"

	"github.com/satori/go.uuid"
)

// _ "github.com/go-sql-driver/mysql" // 仅仅引入, 不使用, 包中的 init() 被执行, 类似Java中的 Class.forname(...), 可以方便的切换不同的驱动

// func openConn() *sql.DB {
// 	conn, err := sql.Open("mysql", "video_user:123456@tcp(localhost:3306)/video_server?charset=utf8")

// 	if err != nil {
// 		panic(err.Error()) // panic 类似Java中的exception, 会中断程序
// 	}

// 	return conn
// }

func AddUserCredential(username string, pwd string) error {
	preparedSql, err := conn.Prepare("insert into user (name, pwd) values (?, ?)")
	if err != nil {
		return err
	}
	_, err = preparedSql.Exec(username, pwd)
	if err != nil {
		return err
	}
	defer preparedSql.Close()
	return nil
}

// find user by username and return pwd
func GetUserCredential(username string) (string, error) {
	preparedSql, err := conn.Prepare("select pwd from user where name = ?")
	if err != nil {
		log.Printf(">>> get user: %s", err)
		return "", err // string 一般没有 nil, 只用 ""表示
	}

	var pwd string
	err = preparedSql.QueryRow(username).Scan(&pwd)
	if err != nil {
		if err == sql.ErrNoRows { // ErrNoRows 表示 查询结果为空, 是特殊的 error
			return "", nil
		}
		return "", err
	}
	defer preparedSql.Close()
	return pwd, nil
}

func DeleteUser(name, pwd string) error {
	preparedSql, err := conn.Prepare("delete from user where name = ? and pwd = ?")
	if err != nil {
		log.Printf(">>> delete user: %s", err)
		return err
	}

	_, err = preparedSql.Exec(name, pwd)
	if err != nil {
		return err
	}
	defer preparedSql.Close()
	return nil
}

///////////////////////////////////////////////

func AddVideo(userId int, videoName string) (*defs.VideoInfo, error) {
	videoId, err := utils.NewUUID()
	if err != nil {
		return nil, err
	}

	time := time.Now()
	timeStr := time.Format("Jan 02 2006, 15:04:05") //固定写法: "Jan 02 2006, 15:04:05"
	preparedSql, err := conn.Prepare(`insert into video_info
        (id, user_id, name, display_ctime) values (?, ?, ?, ?)`)
	if err != nil {
		return nil, err
	}
	_, err = preparedSql.Exec(videoId, userId, videoName, timeStr)
	if err != nil {
		return nil, err
	}
	defer preparedSql.Close()
	return &defs.VideoInfo{Id: videoId, UserId: userId, Name: videoName, DisplayCtime: timeStr}, nil
}

func GetVideo(videoId string) (*defs.VideoInfo, error) {
	preparedSql, err := conn.Prepare("select user_id, name, display_ctime from video_info where id = ?")
	if err != nil {
		return nil, err
	}
	var (
		userId    int
		videoName string
		timeStr   string
	)
	err = preparedSql.QueryRow(videoId).Scan(&userId, &videoName, &timeStr)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	defer preparedSql.Close()
	return &defs.VideoInfo{Id: videoId, UserId: userId, Name: videoName, DisplayCtime: timeStr}, nil
}

func DeleteVideo(videoId string) error {
	preparedSql, err := conn.Prepare("delete from video_info where id = ?")
	if err != nil {
		return err
	}
	_, err = preparedSql.Exec(videoId)
	if err != nil {
		return err
	}
	defer preparedSql.Close()
	return nil
}

/////////////////////////////////////////////////////////////////

func AddComment(videoId string, userId int, content string) (string, error) {
	preparedSql, err := conn.Prepare("insert into comment (id, video_id, user_id, content) values (?, ?, ?, ?)")
	if err != nil {
		return "", err
	}
	commentUUID, err := uuid.NewV4()
	commentId := commentUUID.String()
	if err != nil {
		return "", err
	}

	_, err = preparedSql.Exec(commentId, videoId, userId, content)
	if err != nil {
		return "", err
	}

	defer preparedSql.Close()
	return commentId, nil

}

func ListComment(videoId string, from, to int) ([]*defs.Comment, error) {
	preparedSql, err := conn.Prepare(`select comment.id, user.name, comment.content
		from comment inner join user on comment.user_id = user.id
		where comment.video_id = ?
		and comment.time > FROM_UNIXTIME(?) and comment.time <= FROM_UNIXTIME(?)`)

	var comments []*defs.Comment //return comment list pointer

	rows, err := preparedSql.Query(videoId, from, to)
	if err != nil {
		return comments, err
	}
	for rows.Next() {
		var commentId, username, content string
		if err := rows.Scan(&commentId, &username, &content); err != nil {
			return comments, err
		}
		comment := &defs.Comment{Id: commentId, VideoId: videoId, Username: username, Content: content}
		comments = append(comments, comment)
	}
	return comments, nil
}
