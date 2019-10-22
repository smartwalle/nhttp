package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/smartwalle/http4go"
	"net/http"
)

func main() {

	var s = gin.Default()

	var bufferPool = http4go.NewBufferPool(1024)

	s.Any("/g", func(c *gin.Context) {
		fmt.Println(c.ClientIP(), c.Request.RemoteAddr)
		NewProxy(bufferPool, "127.0.0.1:8891").ServeHTTP(c.Writer, c.Request)
	})

	s.Any("/h", func(c *gin.Context) {
		NewProxy(bufferPool, "127.0.0.1:8892").ServeHTTP(c.Writer, c.Request)
	})

	s.Run(":8899")

}

func NewProxy(bufferPool http4go.BufferPool, to string) *http4go.ReverseProxy {
	var d = func(req *http.Request) {
		req.URL.Scheme = "http"
		req.URL.Host = to
		req.Host = to
		req.Header.Add("X-Real-Ip", req.RemoteAddr)
		fmt.Println(req.URL)
	}

	var p = http4go.NewReverseProxy(bufferPool)
	p.Director = d
	return p
}