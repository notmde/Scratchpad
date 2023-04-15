package api

import (
	"github.com/gin-gonic/gin"
	"github.com/scratchpad-backend/storage"
)

type Server struct {
	listenAddr string
	store      storage.Storage
}

func NewServer(listenAddr string, store storage.Storage) *Server {
	return &Server{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *Server) Run() {
	apiserver := gin.Default()

	gin.SetMode(gin.ReleaseMode)

	apiserver.Use(CORSmanager)

	apiserver.POST("/api/signup", func(c *gin.Context) {
		signup(c, s)
	})
	apiserver.POST("/api/login", func(c *gin.Context) {
		login(c, s)
	})
	apiserver.POST("/api/logout", func(c *gin.Context) {
		VerifyAuth(c, s)
	}, func(c *gin.Context) {
		logout(c, s)
	})

	err := apiserver.Run(s.listenAddr)
	if err != nil {
		panic(err)
	}
}
