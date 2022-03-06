package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func setUpRouter(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
}
