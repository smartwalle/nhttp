package main

import (
	"github.com/gin-gonic/gin"
	"github.com/smartwalle/http4go"
	"net/http"
	"net/url"
)

func main() {
	var s = gin.Default()
	var bufferPool = http4go.NewBufferPool(1024)

	var targets = make(map[string]*url.URL)
	var tURL *url.URL

	tURL, _ = url.Parse("http://127.0.0.1:8891?gate=1")
	targets["user"] = tURL

	tURL, _ = url.Parse("http://127.0.0.1:8892?gate=1")
	targets["order"] = tURL

	var rp = http4go.NewReverseProxy(bufferPool)

	s.Any("/api/:server/*path", func(c *gin.Context) {
		var server = c.Param("server")

		var target = targets[server]
		if target == nil {
			c.AbortWithStatus(http.StatusBadGateway)
			return
		}

		var p = rp.ProxyWithURL(target)
		p.ServeHTTP(c.Writer, c.Request)
	})

	s.Run(":8890")
}
