package server

import (
	"github.com/gin-gonic/gin"
	"github.com/phthaocse/go-gin-demo/middleware"
	"net/http"
)

func setUpRouter(s *Server) {
	r := s.Router
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	userRouter := r.Group("/user")
	{
		userRouter.POST("/register", s.register())
		userRouter.POST("/login", s.login())
		userRouter.GET(":userId", middleware.AuthRequired, s.getUser())
		userRouter.GET("", middleware.AuthRequired, middleware.AdminRequired, s.getAllUser())
	}
}
