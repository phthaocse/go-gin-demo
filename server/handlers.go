package server

import (
	"github.com/gin-gonic/gin"
	"github.com/phthaocse/go-gin-demo/models"
	"github.com/phthaocse/go-gin-demo/schema"
	"github.com/phthaocse/go-gin-demo/utils"
	"net/http"
)

func (s *Server) register() gin.HandlerFunc {
	return func(c *gin.Context) {
		var json schema.UserRegister
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var userid int
		var userEmail string

		err := s.Db.QueryRow(`SELECT email FROM "user" WHERE email = $1`, json.Email).Scan(&userEmail)
		if err == nil {
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
		c.JSON(http.StatusCreated, gin.H{"message": "Register new user successfully"})
	}
}

func (s *Server) login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var userLogin schema.UserLogin
		if err := c.ShouldBindJSON(&userLogin); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		var userEmail string

		err := s.Db.QueryRow(`SELECT email FROM "user" WHERE email = $1`, userLogin.Email).Scan(&userEmail)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Email or Password incorrect"})
			return
		}
	}
}
