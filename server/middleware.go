package server

import (
	"errors"
	"github.com/XiovV/selly-api/jwt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func (s *Server) CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func (s *Server) validateToken() gin.HandlerFunc {

	return func(c *gin.Context) {
		tokenHeader := c.GetHeader("Authorization")

		if len(tokenHeader) == 0 {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "did not receive Authorization header"})
			return
		}

		authorizationHeaderSplit := strings.Split(tokenHeader, " ")
		if len(authorizationHeaderSplit) != 2 && authorizationHeaderSplit[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "wrong Authorization header format"})
			return
		}

		authToken := authorizationHeaderSplit[1]

		token, err := jwt.Validate(authToken)
		if err != nil {
			switch {
			case errors.Is(err, jwt.ErrInvalidSignature):
				s.log.Warnw("invalid signature", "token", authToken)
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "invalid signature"})
			default:
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "invalid token"})
			}
			return
		}

		c.Set("token", token)

		c.Next()
	}
}
