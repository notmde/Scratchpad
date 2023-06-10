package server

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/scratchpad-backend/types"
	"golang.org/x/crypto/bcrypt"
)

func signup(c *gin.Context, s *Server) {
	c.Writer.Header().Set("content-type", "application/json")

	user := new(types.User)

	c.BindJSON(&user)

	if len(user.ID) == 0 {
		c.JSON(400, gin.H{"error": "user id can't be empty"})
		return
	}

	if len(user.Password) == 0 {
		c.JSON(400, gin.H{"error": "password can't be empty"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)

	if err != nil {
		log.Println(err)
		c.JSON(400, gin.H{"error": "failed to hash password"})
		return
	}

	user.Password = string(hashedPassword)

	if err = s.store.InsertDB(user); err != nil {
		log.Println(err)
		c.JSON(400, gin.H{"error": "user already exists"})
		return
	}

	c.JSON(201, gin.H{"status": "user created"})
}

func login(c *gin.Context, s *Server) {
	c.Writer.Header().Set("content-type", "application/json")

	body := new(types.User)

	c.BindJSON(&body)

	if len(body.ID) == 0 {
		c.JSON(400, gin.H{"error": "user id can't be empty"})
		return
	}

	if len(body.Password) == 0 {
		c.JSON(400, gin.H{"error": "password can't be empty"})
		return
	}

	user := new(types.User)

	if err := s.store.GetDB(body.ID, user); err != nil {
		c.JSON(400, gin.H{"error": "user not found"})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		c.JSON(401, gin.H{"error": "invalid password"})
		return
	}

	// jwt token generation
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  body.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))

	if err != nil {
		c.JSON(500, gin.H{"error": "cannot create tokenstring"})
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "authentication",
		Value:    tokenString,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})

	c.JSON(200, gin.H{"status": "logged in", "canvas_data": user.CanvasData})
}

func logout(c *gin.Context, s *Server) {
	c.Writer.Header().Set("content-type", "application/json")

	if !VerifyAuth(c, s) {
		return
	}

	body := new(types.User)

	c.BindJSON(&body)

	err := s.store.UpdateCanvasDB(body.ID, body.CanvasData)

	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{"error": "failed to update data"})
		return
	}

	// delete cookie
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "authentication",
		Value:    "Logout",
		Expires:  time.Now(),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})

	c.JSON(200, gin.H{"status": "logged out"})
}
