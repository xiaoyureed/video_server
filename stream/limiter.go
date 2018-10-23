// 流控：bucket token 算法
package main

import (
	"log"
)

type ConnLimiter struct {
	concurrentConn int
	bucket         chan int
}

// NewConnLimiter constructor
func NewConnLimiter(cc int) *ConnLimiter {
	return &ConnLimiter{
		concurrentConn: cc,
		bucket:         make(chan int, cc),
	}
}

// GetConn get stream conn
func (cl *ConnLimiter) GetConn() bool {
	if len(cl.bucket) >= cl.concurrentConn {
		log.Printf("Reached the limit")
		return false
	}

	cl.bucket <- 1
	log.Printf("Get conn")
	return true
}

//ReleaseConn give back stream conn
func (cl *ConnLimiter) ReleaseConn() {
	c := <-cl.bucket
	log.Printf("Release conn : %d", c)
}
