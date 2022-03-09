package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/phthaocse/go-gin-demo/schema"
	"net/http"
)

func (s *Server) register() gin.HandlerFunc {
	return func(c *gin.Context) {
		var json schema.Register
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var userid int
		res, _ := s.Db.Exec(`SELECT email FROM "user" WHERE email = $1`, json.Email)
		if res != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "User has been existed"})
			return
		}
		err := s.Db.QueryRow(`INSERT INTO "user" (username, email)
									VALUES($1, $2) RETURNING id`, json.Username, json.Email).Scan(&userid)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		fmt.Println(userid)
	}
}
