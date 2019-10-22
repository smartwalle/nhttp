package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	var s = gin.Default()
	s.GET("/api/user/list", func(c *gin.Context) {
		c.String(200, "user list")
	})
	s.GET("/api/user/detail", func(c *gin.Context) {
		c.String(200, "user detail")
	})
	s.Run(":8891")
}
