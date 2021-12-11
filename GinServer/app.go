package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const TEMPLATE_DIR = "./templates/"

func main() {
	r := gin.Default()

	// Quick Stark
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Response Web Page
	r.LoadHTMLGlob(TEMPLATE_DIR + "*.html")
	r.GET("/", func(ctx *gin.Context) {
		ctx.HTML(200, "index.html", gin.H{})
	})
	r.GET("/xss", func(ctx *gin.Context) {
		ctx.HTML(200, "xss.html", gin.H{})
	})
	r.POST("/xss", func(ctx *gin.Context) {
		data := ctx.PostForm("data")
		// そのまま表示してもデフォルトでHTMLタグはサニタイズされた
		ctx.HTML(200, "xss.html", gin.H{"data": data})
	})

	//parameter
	r.GET("/user/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.String(http.StatusOK, "Hello %s", name)
	})
	r.GET("/user/:name/*action", func(c *gin.Context) {
		name := c.Param("name")
		action := c.Param("action")
		message := name + " is " + action
		c.String(http.StatusOK, message)
	})

	//default port: 8000
	r.Run(":8000")
}
