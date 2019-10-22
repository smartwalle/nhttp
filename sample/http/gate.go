package main

import (
	"github.com/gin-gonic/gin"
	"github.com/smartwalle/http4go"
	"net/http"
)

func main() {
	var proxyList = make(map[string]*ProxyInfo)
	proxyList["user"] = &ProxyInfo{Scheme: "http", Host: "127.0.0.1:8891"}
	proxyList["order"] = &ProxyInfo{Scheme: "http", Host: "127.0.0.1:8892"}

	var s = gin.Default()
	var bufferPool = http4go.NewBufferPool(1024)

	s.Any("/api/:server/*path", func(c *gin.Context) {
		var server = c.Param("server")
		var path = c.Param("path")

		var proxy = proxyList[server]

		var toPath = "/api/" + server + path
		NewProxy(bufferPool, proxy.Scheme, proxy.Host, toPath).ServeHTTP(c.Writer, c.Request)
	})

	s.Run(":8890")
}

type ProxyInfo struct {
	Scheme string
	Host   string
}

func NewProxy(bufferPool http4go.BufferPool, scheme, host, path string) *http4go.ReverseProxy {
	var d = func(req *http.Request) {
		req.URL.Scheme = scheme
		req.URL.Host = host
		req.URL.Path = path
	}

	var p = http4go.NewReverseProxy(bufferPool)
	p.Director = d
	return p
}
