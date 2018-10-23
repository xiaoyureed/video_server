package defs

type UserCredential struct {
	Username string `json:"user_name"` // go 内置json序列化
	Pwd      string `json:"pwd"`
}

type Signedup struct {
	Success   bool   `json:"success"`
	SessionId string `json:"session_id"`
}

// video info model
type VideoInfo struct {
	Id           string
	UserId       int
	Name         string
	DisplayCtime string
}

type Comment struct {
	Id       string
	VideoId  string
	Username string
	Content  string
}

type SimpleSession struct {
	Username string // login name
	Ttl      int64
}
