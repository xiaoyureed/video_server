package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"video_server/api/dbops"
	"video_server/api/defs"
	"video_server/api/session"

	"github.com/julienschmidt/httprouter"
)

// CreateUser create user
func CreateUser(writer http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// io.WriteString(writer, ">>> create user")

	res, _ := ioutil.ReadAll(req.Body)
	user := &defs.UserCredential{}
	if err := json.Unmarshal(res, user); err != nil {
		sendErrorResponse(writer, defs.ErrorRequestBodyParseFailed)
		return
	}

	if err := dbops.AddUserCredential(user.Username, user.Pwd); err != nil {
		sendErrorResponse(writer, defs.ErrorDBOperationFailed)
		return
	}

	sessionID := session.GenerateSessionId(user.Username)
	signup := &defs.Signedup{SessionId: sessionID, Success: true}
	if respByteArr, err := json.Marshal(signup); err != nil {
		sendErrorResponse(writer, defs.ErrorInternalError)
	} else {

		sendNormalResponse(writer, string(respByteArr), 201)
	}
}

// Login login
func Login(writer http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// username := params.ByName("user_name")
	// io.WriteString(writer, username)

	reqData, _ := ioutil.ReadAll(req.Body)
	user := &defs.UserCredential{}

	if err := json.Unmarshal(reqData, user); err != nil {
		sendErrorResponse(writer, defs.ErrorRequestBodyParseFailed)
		return
	}
	pwd, err := dbops.GetUserCredential(user.Username)
	if err != nil {
		sendErrorResponse(writer, defs.ErrorDBOperationFailed)
		return
	}
	if user.Pwd != pwd { // wrong pwd
		sendErrorResponse(writer, defs.ErrorNotAuthUser)
		return
	}

	// create session
	sessionID := session.GenerateSessionId(user.Username)
	sendNormalResponse(writer, sessionID, 201)
}
