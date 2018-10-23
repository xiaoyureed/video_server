package dbops

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

// clear data at the begining of the test
func clearTables() {
	conn.Exec("truncate user")
	conn.Exec("truncate video_info")
	conn.Exec("truncate comment")
	conn.Exec("truncate session")
}

func TestMain(m *testing.M) {
	clearTables()
	m.Run()
	clearTables()
}

func TestUserWorkflow(t *testing.T) {
	clearTables()
	t.Run("add", testAddUser)
	t.Run("get", testGetUser)
	t.Run("delete", testDeleteUser)
}

func TestVideoWorkflow(t *testing.T) {
	clearTables()
	t.Run("add", testAddVideo)
	t.Run("get", testGetVideo)
	t.Run("delete", testDeleteVideo)
}

func TestCommentWorkflow(t *testing.T) {
	clearTables()
	t.Run("add", testAddComment)
	t.Run("list", testListComment)
}

func TestSessionWorkflow(t *testing.T) {
	clearTables()
	t.Run("insert", testInsertSession)
	t.Run("retrieve", testRetrieveSession)
	t.Run("retrieve all", testRetrieveAllSession)
	t.Run("delete", testDeleteSession)

}

func testInsertSession(t *testing.T) {
	sessionId := "digng-sidng-sifng"
	ttl := int64(123940384)
	username := "xiaoyu"

	err := InsertSession(sessionId, ttl, username)
	if err != nil {
		t.Errorf("Error of InsertSession: %v", err)
	}
}

func testRetrieveSession(t *testing.T) {
	sessionId := "digng-sidng-sifng"
	ttl := int64(123940384)
	username := "xiaoyu"

	session, err := RetrieveSession(sessionId)
	if err != nil {
		t.Errorf("Error of RetrieveSession: %v", err)
	}
	if session.Ttl != ttl || session.Username != username {
		t.Errorf("RetrieveSession failed")
	}
}

func testRetrieveAllSession(t *testing.T) {
	sessionId := "digng-sidng-sifng"
	// ttl := 123940384
	// username := "xiaoyu"

	sessionMap, err := RetrieveAllSession()
	if err != nil {
		t.Errorf("Error of RetrieveSession: %v", err)
	}
	sessionMap.Range(func(key, value interface{}) bool {
		// session := value.(*defs.SimpleSession)
		if key != sessionId {
			t.Errorf("RetrieveAllSession failed")
		}
		return true
	})
}

func testDeleteSession(t *testing.T) {
	sessionId := "digng-sidng-sifng"
	err := DeleteSession(sessionId)
	if err != nil {
		t.Errorf("Error of DeleteSession: %v", err)
	}

	session, err := RetrieveSession(sessionId)
	if session != nil {
		t.Errorf("DeleteSession failed")
	}
}

func testAddComment(t *testing.T) {
	videoId := "sidnglwig"
	userId := 1
	content := "I like it."

	commentId, err := AddComment(videoId, userId, content)
	if err != nil {
		t.Errorf("Error of AddComment: %v", err)
	}
	if commentId == "" {
		t.Errorf("AddComment failed")
	}
}

func testListComment(t *testing.T) {
	videoId := "sidnglwig"
	from := 1514764800
	to, _ := strconv.Atoi(strconv.FormatInt(time.Now().UnixNano()/1000000000, 10))

	resp, err := ListComment(videoId, from, to)
	if err != nil {
		t.Errorf("Error of ListComment: %v", err)
	}
	for index, element := range resp {
		fmt.Printf(">>> comment: %d, %v \n", index, element)
	}
}

var tempVideoId string

func testAddVideo(t *testing.T) {
	video, err := AddVideo(1, "1-video")
	if err != nil {
		t.Errorf("Error of AddVideo: %v", err)
	}
	tempVideoId = video.Id
}

func testGetVideo(t *testing.T) {
	video, err := GetVideo(tempVideoId)
	if err != nil {
		t.Errorf("Error of GetVideo: %v", err)
	}
	if video.UserId != 1 || video.Name != "1-video" {
		t.Errorf("Bussiness error of GetVideo")
	}
}

func testDeleteVideo(t *testing.T) {
	err := DeleteVideo(tempVideoId)
	if err != nil {
		t.Errorf("Error of DeleteVideo: %v", err)
	}

	video, _ := GetVideo(tempVideoId)
	if video != nil {
		t.Errorf("delete video failed")
	}
}

func testAddUser(t *testing.T) {
	err := AddUserCredential("xiaoyu", "123")
	if err != nil {
		t.Errorf("Error of AddUser: %v", err)
	}
}

func testGetUser(t *testing.T) {
	pwd, err := GetUserCredential("xiaoyu")
	if pwd != "123" || err != nil {
		t.Errorf("Error of GetUser: %v", err)
	}
}

func testDeleteUser(t *testing.T) {
	err := DeleteUser("xiaoyu", "123")
	if err != nil {
		t.Errorf("Error of DeleteUser: %v", err)
	}

	pwd, _ := GetUserCredential("xiaoyu")
	if pwd != "" {
		t.Errorf("DeleteUser failed")
	}
}
