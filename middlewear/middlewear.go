package middlewear

import (
	"net/http"

	st "VK_app/structures"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// CheckToken is a function that takes an http.Handler and returns an http.Handler.
//
// It checks the validity of the token in the request header and calls the next handler if the token is valid.

func CheckToken(c *gin.Context) {
	tokenString := c.Request.Header.Get("Authorization")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return st.Secret, nil
	})
	if err != nil || !token.Valid {
		log.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized user"})
		return
	}
	c.Next()
}

// CheckTokenAdmin checks the token in the Gin context for admin access.
//
// Parameter(s):
//
//	c *gin.Context
func CheckTokenAdmin(c *gin.Context) {
	tokenString := c.Request.Header.Get("Authorization")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return st.Secret, nil
	})
	if err != nil || !token.Valid {
		log.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized admin"})
		return
	}
	if role, ok := token.Claims.(jwt.MapClaims)["role"].(int); ok && role != 1 {
		log.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized admin role"})
		return
	}
}
