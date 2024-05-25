package hello

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func init() {
	r := gin.New()

	r.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Hello world!")
	})

	r.GET("/ping", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "pong")
	})

	http.Handle("/", r)
}
