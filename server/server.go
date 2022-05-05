package server

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Server struct {
	log *zap.SugaredLogger
}

func New(log *zap.SugaredLogger) *Server {
	return &Server{log: log}
}

func (s *Server) Serve() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery(), s.CORS())

	v1 := router.Group("/v1")

	usersPublic := v1.Group("/users")
	{
		usersPublic.GET("/token", s.generateToken)
		usersPublic.GET("/refresh-token", s.refreshToken)
	}

	return router
}
