package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func StartServer() {
	log.Println("Start server")

	r := gin.Default()

	r.LoadHTMLGlob("templates/*")

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message" : "pong",
		})
	})

	r.GET("/hello", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", gin.H{})
	})
	
	r.Run()

	log.Println("Server down")
}