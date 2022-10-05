package server

import (
	"errors"
	"github.com/XiovV/selly-api/jwt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func (s *Server) getMissedMessages(c *gin.Context) {
	sellyId := s.getContextSellyID(c)

	messages := s.redis.GetMessages(sellyId)

	go s.redis.DelUser(sellyId)

	c.JSON(http.StatusOK, gin.H{"messages": messages})
}

func (s *Server) generateToken(c *gin.Context) {
	hashedSeed := c.Query("id")

	if len(hashedSeed) < 64 {
		s.log.Warnw("selly id query parameter has an invalid length", "expected", 64, "got", len(hashedSeed), "id", hashedSeed)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	sellyId := s.generateSellyID(hashedSeed)

	newToken, err := jwt.New(sellyId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": newToken})
}

func (s *Server) refreshToken(c *gin.Context) {
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
		if errors.Is(err, jwt.ErrInvalidSignature) {
			s.log.Warnw("invalid signature", "token", authToken)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "invalid signature"})
			return
		}
	}

	if !jwt.IsTokenExpired(token) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token still hasn't expired"})
		return
	}

	sellyId := jwt.GetClaimString(token, "sellyId")

	newToken, err := jwt.New(sellyId)
	if err != nil {
		s.log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": newToken})
}
