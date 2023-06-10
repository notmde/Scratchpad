package server

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func VerifyAuth(c *gin.Context, s *Server) {
	tokenString, err := c.Cookie("authentication")

	if err != nil {
		c.JSON(401, gin.H{"error": "no token"})
		return
	}

	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			c.JSON(500, gin.H{"error": "unexpected signing method"})
		}

		return []byte(os.Getenv("SECRET_KEY")), nil
		//return []byte(SECRET_KEY), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.JSON(401, gin.H{"error": "token expired"})
			return
		}

		err := s.store.IfIDExistsDB(claims["id"].(string))

		if err != nil {
			c.JSON(404, gin.H{"error": "id not found"})
			return
		}

	} else {
		c.JSON(401, gin.H{"error": "invalid token"})
		return
	}

	c.Next()
}

func CORSmanager(c *gin.Context) {
	origin := c.GetHeader("Origin")

	c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204)
		return
	}

	c.Next()
}
