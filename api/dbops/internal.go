// session db 操作相关
package dbops

import (
	"database/sql"
	"log"
	"strconv"
	"sync"
	"video_server/api/defs"
)

func InsertSession(sessionId string, ttl int64, username string) error {
	ttlStr := strconv.FormatInt(ttl, 10)
	preparedSql, err := conn.Prepare("insert into session (id, ttl, login_name) values (?, ?, ?)")

	_, err = preparedSql.Exec(sessionId, ttlStr, username)
	if err != nil {
		return err
	}
	defer preparedSql.Close()
	return nil
}

func DeleteSession(sessionId string) error {
	preparedSql, err := conn.Prepare("delete from session where id = ?")
	_, err = preparedSql.Exec(sessionId)
	if err != nil {
		return err
	}
	return nil
}

func RetrieveSession(sessionId string) (*defs.SimpleSession, error) {
	session := &defs.SimpleSession{} // 这个不是类型, 直接是一个变量了, 推荐这种方式
	preparedSql, err := conn.Prepare("select ttl, login_name from session where id = ?")
	var ttlStr, username string
	err = preparedSql.QueryRow(sessionId).Scan(&ttlStr, &username)
	if err != nil && err != sql.ErrNoRows {
		return session, err
	}
	if ttl, err := strconv.ParseInt(ttlStr, 10, 64); err == nil {
		session.Ttl = ttl
		session.Username = username
	} else {
		return nil, err
	}
	defer preparedSql.Close()
	return session, nil
}

func RetrieveAllSession() (*sync.Map, error) {
	sessionMap := &sync.Map{}
	preparedSql, err := conn.Prepare("select id, ttl, login_name from session")

	rows, err := preparedSql.Query()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var id, ttlStr, username string
		if err := rows.Scan(&id, &ttlStr, &username); err != nil {
			log.Printf("retrive all session error: %s", err)
			break
		}
		if ttl, err := strconv.ParseInt(ttlStr, 10, 64); err == nil {
			session := &defs.SimpleSession{Username: username, Ttl: ttl}
			sessionMap.Store(id, session)
			log.Printf("session id: %s, ttl: %d", id, session.Ttl)
		}

	}
	return sessionMap, nil

}
