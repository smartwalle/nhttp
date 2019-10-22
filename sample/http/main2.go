package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	var s = gin.Default()
	s.GET("g", func(c *gin.Context) {
		c.String(200, "8891")
		fmt.Println(c.ClientIP(), c.Request.RemoteAddr)
		fmt.Println(c.Request.Header)
	})
	s.Run(":8891")
}