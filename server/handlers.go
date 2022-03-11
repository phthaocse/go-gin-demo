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
		var userReg schema.UserRegister
		if err := c.ShouldBindJSON(&userReg); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user := models.User{Email: userReg.Email, Password: userReg.Password, Username: userReg.Username}
		if user.IsExist(s.Db) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "User has been existed"})
			return
		}

		_, err := user.Create(s.Db)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "Register new user successfully"})
		return
	}
}

func (s *Server) login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var userLogin schema.UserLogin
		if err := c.ShouldBindJSON(&userLogin); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		user := models.User{Email: userLogin.Email, Password: userLogin.Password}
		if !user.IsExist(s.Db) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Email or Password incorrect"})
			return
		}
		if !utils.CheckPasswordHash(userLogin.Password, user.Password) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Email or Password incorrect"})
			return
		}
		jwt := utils.GetJWT(s.Config.SecretKey, user.Id)
		c.JSON(http.StatusOK, gin.H{"access_token": jwt})
	}
}

func (s *Server) getUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		param := struct {
			UserId int `uri:"userId" binding:"required"`
		}{}
		if err := c.ShouldBindUri(&param); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		reqUser, ok := c.Get("CurrUser")
		if !ok || reqUser != param.UserId {
			c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission on this resource"})
			return
		}
		user := models.User{Id: param.UserId}
		err := user.GetByPk(s.Db)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		c.JSON(http.StatusOK, gin.H{"user": user})
	}
}
