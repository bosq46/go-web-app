package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("../templates/*.html")

	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(200, "index.html", gin.H{})
	})

	router.GET("/xss", func(ctx *gin.Context) {
		ctx.HTML(200, "xss.html", gin.H{})
	})
	router.POST("/xss", func(ctx *gin.Context) {
		data := ctx.PostForm("data")
		// そのまま表示してもデフォルトでHTMLタグはサニタイズされた
		ctx.HTML(200, "xss.html", gin.H{"data": data})
	})
	router.Run()
}
