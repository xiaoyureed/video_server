package main

import (
	"encoding/json"
	"io"
	"net/http"
	"video_server/api/defs"
)

// send error resp
func sendErrorResponse(writer http.ResponseWriter, resp defs.ErrorResponse) {
	writer.WriteHeader(resp.HttpCode)
	respByte, _ := json.Marshal(resp.Error)
	io.WriteString(writer, string(respByte))
}

// send normal resp
func sendNormalResponse(writer http.ResponseWriter, resp string, httpCode int) {
	writer.WriteHeader(httpCode)
	io.WriteString(writer, resp)
}
