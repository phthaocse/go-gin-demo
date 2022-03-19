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
			c.JSON(http.StatusBadRequest, gin.H{"error": "User has been existed"})
			return
		}
		if user.IsUsernameExisted(s.Db) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username has been existed"})
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
		jwt := utils.GetJWT(s.Config.SecretKey, user.Id, user.Role)
		c.JSON(http.StatusOK, gin.H{"access_token": jwt})
	}
}

func (s *Server) getUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		param := struct {
			UserId int `uri:"userId" binding:"required"`
		}{}
		if err := c.ShouldBindUri(&param); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

func (s *Server) getAllUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var query schema.DefaultQuery
		if err := c.ShouldBindQuery(&query); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		user := &models.User{}
		var err error
		var users []*models.User
		if query.Limit > 0 {
			users, err = user.GetMulti(s.Db, query.Limit, query.Offset)
		} else {
			users, err = user.GetAll(s.Db)
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		c.JSON(http.StatusOK, gin.H{"users": users})
		return
	}
}

func (s *Server) UpdateActiveStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		param := struct {
			UserId int `uri:"userId" binding:"required"`
		}{}
		if err := c.ShouldBindUri(&param); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var payload = struct {
			ActiveStatus *bool `json:"active_status" binding:"required"`
		}{}
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		user := &models.User{Id: param.UserId}
		err := user.UpdateActiveStatus(s.Db, *payload.ActiveStatus)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"user": user})
		return
	}
}

func (s *Server) createTicket() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ticketPayload schema.Ticket
		if err := c.ShouldBindJSON(&ticketPayload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ticket := &models.Ticket{Title: ticketPayload.Title, Assignee: ticketPayload.Assignee, Content: ticketPayload.Content}

		if val, ok := c.Get("CurrUser"); ok {
			ticket.Reporter = val.(int)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Can't get request user"})
		}

		err := ticket.Create(s.Db)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"data": ticket})
		return
	}
}
