package main

import (
	"example/chatapp-go-api/pkg/websocket"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}

	router.Use(cors.New(config))

	pool := websocket.NewPool()
	go pool.Start()
	router.GET("/ws/:username", func(c *gin.Context) {
		username := c.Param("username")
		wshandler(c.Writer, c.Request, pool, username)
	})
	router.GET("/", HomeHandler)
	router.GET("/users", GetUsers)
	router.POST("/users", AddUser)
	router.DELETE("/users/:username", DeleteUserByUsername)
	router.POST("/login", Login)
	router.Run(":8080")
}
