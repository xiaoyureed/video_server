package session

import (
	"sync"
	"time"
	"video_server/api/dbops"
	"video_server/api/defs"
	"video_server/api/utils"
)

var sessionMap *sync.Map

func init() {
	sessionMap = &sync.Map{}
}

func LoadSessionFromDB() {
	m, err := dbops.RetrieveAllSession()
	if err != nil {
		return
	}

	// m.Range(func(key, value interface{}) bool {
	// 	session := value.(*defs.SimpleSession)
	// 	sessionMap.Store(key, session)
	// 	return true
	// })

	sessionMap = m
}

// create session and return session id
func GenerateSessionId(username string) string {
	sessionId, _ := utils.NewUUID()
	createTime := nowMilliSecond()
	ttl := createTime + 30*60*1000 // serverside session valide time: 30min

	dbops.InsertSession(sessionId, ttl, username)
	sessionMap.Store(sessionId, &defs.SimpleSession{Ttl: ttl, Username: username})
	return sessionId

}

func IsSessionExpired(sessionId string) (string, bool) {
	// session, err := dbops.RetrieveSession(sessionId)

	sessionRaw, ok := sessionMap.Load(sessionId)
	if !ok {
		return "", true
	}

	session := sessionRaw.(*defs.SimpleSession)
	if session.Ttl < nowMilliSecond() { // session expired
		deleteExpiredSession(sessionId)
		return "", true
	}
	return session.Username, false
}

func deleteExpiredSession(sessionId string) {
	sessionMap.Delete(sessionId)
	dbops.DeleteSession(sessionId)
}

func nowMilliSecond() int64 {
	return time.Now().UnixNano() / 1000000 //unit: ms
}
