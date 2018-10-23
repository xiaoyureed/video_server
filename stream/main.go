package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type middleWareHandler struct {
	router  *httprouter.Router
	limiter *ConnLimiter
}

func newMiddleWareHandler(r *httprouter.Router, cc int) http.Handler {
	return middleWareHandler{
		limiter: NewConnLimiter(cc),
		router:  r,
	}
	// m := middleWareHandler{}
	// m.limiter = NewConnLimiter(cc)
	// m.router = r
	// return m
}

func (m middleWareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !m.limiter.GetConn() {
		sendErrorResponse(w, http.StatusTooManyRequests, "Too many request.")
		return
	}
	m.router.ServeHTTP(w, r)
	defer m.limiter.ReleaseConn()
}

func registerHandler() *httprouter.Router {
	router := httprouter.New()
	router.GET("/videos/:id", streamHandler)
	router.POST("/videos/:id", uploadHandler)
	router.GET("/upload-page", uploadPageHandler)

	return router
}

func main() {
	router := registerHandler()
	mw := newMiddleWareHandler(router, 2)
	http.ListenAndServe(":9000", mw)
}
