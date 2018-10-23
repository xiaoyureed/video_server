package main

import (
	"net/http"
	"video_server/api/defs"
	"video_server/api/session"
)

var HEADER_FIELD_SESSION = "X-session-id"
var HEADER_FIELD_UNAME = "X-user-name"

func validateSession(req *http.Request) bool {
	sessionId := req.Header.Get(HEADER_FIELD_SESSION)
	if len(sessionId) == 0 {
		return false
	}
	uname, ok := session.IsSessionExpired(sessionId)
	if ok {
		return false
	}
	req.Header.Set(HEADER_FIELD_UNAME, uname) // session validateion success
	return true
}

func validateUser(w *http.ResponseWriter, req *http.Request) bool {
	uname := req.Header.Get(HEADER_FIELD_UNAME)
	if len(uname) == 0 {
		sendErrorResponse(*w, defs.ErrorNotAuthUser)
		return false
	}
	return true
}
