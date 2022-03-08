package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func setUpRouter(s *server) {
	r := s.router
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.POST("/register", s.register())
}
