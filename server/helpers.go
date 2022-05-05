package server

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func (s *Server) getContextToken(c *gin.Context) *jwt.Token {
	t, exists := c.Get("token")

	if !exists {
		s.log.Error("token not found in context")
		return nil
	}

	token := t.(*jwt.Token)

	return token
}
