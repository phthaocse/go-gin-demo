package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/phthaocse/go-gin-demo/utils"
	"net/http"
	"strings"
)

func AuthRequired(c *gin.Context) {
	var authHeader = struct {
		Authorization string
	}{}
	if err := c.ShouldBindHeader(&authHeader); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Could not validate credential"})
		return
	}
	tokenString := strings.Split(authHeader.Authorization, "Bearer ")
	if len(tokenString) == 2 {
		if claim, err := utils.ParseJWT(tokenString[1], []byte(utils.GetEnv("SECRET_KEY", ""))); err != nil {
			fmt.Println(err.Error())
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Could not validate credential"})
			return
		} else {
			c.Set("CurrUser", claim.UserId)
			return
		}
	}

	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Could not validate credential"})
	return
}
