package main

import (
	"github.com/gin-gonic/gin"
)

func main() {

	go router()

	r := gin.New()
	r.GET("/ws", ginWsServe())
	r.Run(":3000")
}

func ginWsServe() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		handlesWebSocketRequests(c.Writer, c.Request)
	})
}