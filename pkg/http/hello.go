package httpX

import "github.com/gin-gonic/gin"

func Hello(c *gin.Context) {
	ip := c.ClientIP()
	c.Writer.WriteString("hello " + ip)
}
