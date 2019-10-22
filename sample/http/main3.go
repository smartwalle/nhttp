package main

import "github.com/gin-gonic/gin"

func main() {
	var s = gin.Default()
	s.GET("h", func(c *gin.Context) {
		c.String(200, "8892")
	})
	s.Run(":8892")
}
