package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// 中间层, 检查 session有效, user有效...
type middlewareHandler struct {
	router *httprouter.Router
}

func (m middlewareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// check session
	validateSession(r)

	m.router.ServeHTTP(w, r)
}

func createMiddlewareHandler(router *httprouter.Router) http.Handler {
	m := middlewareHandler{}
	m.router = router
	return m
}

// register handler
func RegisterHandler() *httprouter.Router {
	router := httprouter.New()
	router.POST("/users", CreateUser) // register
	router.POST("/logins", Login)     // login
	return router
}

func main() {
	r := RegisterHandler()
	handler := createMiddlewareHandler(r)
	http.ListenAndServe(":8080", handler)
}
