package dbops

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql" // 仅仅引入, 不使用, 包中的 init() 被执行, 类似Java中的 Class.forname(...), 可以方便的切换不同的驱动
)

var (
	conn *sql.DB
	err  error
)

func init() {
	conn, err = sql.Open("mysql", "video_user:123456@tcp(localhost:3306)/video_server?charset=utf8")

	if err != nil {
		panic(err.Error()) // panic 类似Java中的exception, 会中断程序
	}
}
