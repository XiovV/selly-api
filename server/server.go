package server

import (
	"github.com/XiovV/selly-api/redis"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Server struct {
	redis *redis.Redis
	log   *zap.SugaredLogger
}

func New(redis *redis.Redis, log *zap.SugaredLogger) *Server {
	return &Server{redis: redis, log: log}
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

	usersProtected := v1.Group("/users")
	usersProtected.Use(s.validateToken())
	{
		usersProtected.GET("/missed-messages", s.getMissedMessages)
	}

	return router
}
