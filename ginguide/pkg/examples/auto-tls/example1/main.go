package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/autotls"

)

func main() {
	r := gin.Default()

	r.GET("/ping", func(ctx *gin.Context) {
		ctx.String(200, "pong")
	})

	log.Fatal(autotls.Run(r, "example1.com", "example2.com"))
}
