package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/smartwalle/nhttp"
	"net/http"
	"net/url"
)

func main() {
	var s = gin.Default()
	var bufferPool = nhttp.NewBufferPool(1024)

	var targets = make(map[string]*url.URL)
	var tURL *url.URL

	tURL, _ = url.Parse("http://127.0.0.1:8891?gate=1")
	targets["user"] = tURL

	tURL, _ = url.Parse("http://127.0.0.1:8892?gate=1")
	targets["order"] = tURL

	var rp = nhttp.NewReverseProxy(bufferPool, nil)

	s.Any("/api/:server/*path", func(c *gin.Context) {
		var server = c.Param("server")

		var target = targets[server]
		if target == nil {
			c.AbortWithStatus(http.StatusBadGateway)
			return
		}

		// 如果有需要解析参数，则需要把 request body 复制一份
		nBody, err := nhttp.DumpBody(c.Request)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer func() {
			// 重新设置 request body
			c.Request.Body = nBody
		}()

		// 解析参数
		c.Request.ParseForm()
		fmt.Println("p1", "-", c.Request.Form["p1"])
		fmt.Println("p2", "-", c.Request.Form["p2"])

		var p = rp.ProxyWithURL(target)
		p.ServeHTTP(c.Writer, c.Request)
	})

	s.Run(":8890")
}
