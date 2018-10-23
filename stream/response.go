package main

import (
	"io"
	"net/http"
)

func sendErrorResponse(w http.ResponseWriter, httpCode int, msg string) {
	w.WriteHeader(httpCode)
	io.WriteString(w, msg)
}
