package main

import "github.com/gin-gonic/gin"

func main() {
	var s = gin.Default()
	s.GET("/api/order/list", func(c *gin.Context) {
		c.String(200, "order list")
	})
	s.GET("/api/order/detail", func(c *gin.Context) {
		c.String(200, "order detail")
	})
	s.Run(":8892")
}
