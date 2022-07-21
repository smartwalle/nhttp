package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	var s = gin.Default()
	s.GET("/api/order/list", func(c *gin.Context) {
		c.Request.ParseForm()
		fmt.Println("接口:", c.Request.URL.Path)
		fmt.Println("收到请求参数:")
		for key := range c.Request.Form {
			fmt.Println(key, "-", c.Request.Form.Get(key))
		}
		c.String(200, "order list")
	})
	s.GET("/api/order/detail", func(c *gin.Context) {
		c.Request.ParseForm()
		fmt.Println("接口:", c.Request.URL.Path)
		fmt.Println("收到请求参数:")
		for key := range c.Request.Form {
			fmt.Println(key, "-", c.Request.Form.Get(key))
		}
		c.String(200, "order detail")
	})
	s.Run(":8892")
}
