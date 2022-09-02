package server

import (
	"crypto/sha256"
	"fmt"
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

func (s *Server) generateSellyID(hashedSeed string) string {
	sellyId := sha256.Sum256([]byte(hashedSeed))

	return fmt.Sprintf("%x", sellyId[:])
}
