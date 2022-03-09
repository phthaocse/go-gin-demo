package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/phthaocse/go-gin-demo/models"
	"github.com/phthaocse/go-gin-demo/schema"
	"github.com/phthaocse/go-gin-demo/utils"
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
		var userEmail string
		err := s.Db.QueryRow(`SELECT email FROM "user" WHERE email = $1`, json.Email).Scan(&userEmail)
		if userEmail != "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "User has been existed"})
			return
		}
		hashPwd, err := utils.HashPassword(json.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "err"})
		}
		err = s.Db.QueryRow(`INSERT INTO "user" (username, email, password, role)
									VALUES($1, $2, $3, $4) RETURNING id`,
			json.Username, json.Email, hashPwd, models.MemberRole).Scan(&userid)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		fmt.Println(userid)
	}
}
