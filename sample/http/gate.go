package main

import (
	"github.com/gin-gonic/gin"
	"github.com/smartwalle/http4go"
	"net/url"
)

func main() {
	var s = gin.Default()
	var bufferPool = http4go.NewBufferPool(1024)

	var rp = http4go.NewReverseProxy(bufferPool)

	s.Any("/api/user/*path", func(c *gin.Context) {
		var u, _ = url.Parse("http://127.0.0.1:8891?c=2")
		rp.ProxyWithURL(u).ServeHTTP(c.Writer, c.Request)
	})

	s.Any("/api/order/*path", func(c *gin.Context) {
		var u, _ = url.Parse("http://127.0.0.1:8892?c=2")
		rp.ProxyWithURL(u).ServeHTTP(c.Writer, c.Request)
	})

	s.Run(":8890")
}
