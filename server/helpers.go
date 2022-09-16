package server

import (
	"crypto/sha256"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
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

func (s *Server) getContextSellyID(c *gin.Context) string {
	sellyId, exists := c.Get("sellyId")
	if !exists {
		s.log.Errorw("selly id not in context")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return ""
	}

	return sellyId.(string)
}

func (s *Server) generateSellyID(hashedSeed string) string {
	sellyId := sha256.Sum256([]byte(hashedSeed))

	return fmt.Sprintf("%x", sellyId[:])
}
