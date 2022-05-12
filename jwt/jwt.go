package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"os"
	"strings"
	"time"
)

func New(sellyId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sellyId": sellyId,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	})

	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func Validate(tok string) (*jwt.Token, error) {
	token, err := jwt.Parse(tok, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("couldn't verify signing method")
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		switch {
		case strings.Contains(err.Error(), "signature is invalid"):
			return nil, ErrInvalidSignature
		case strings.Contains(err.Error(), "Token is expired"):
			return token, ErrTokenExpired
		default:
			return nil, err
		}
	}

	return token, nil
}

func IsTokenExpired(token *jwt.Token) bool {
	return getExp(token) < time.Now().Unix()
}

func getExp(token *jwt.Token) int64 {
	return int64(token.Claims.(jwt.MapClaims)["exp"].(float64))
}

func GetClaimString(token *jwt.Token, claim string) string {
	return token.Claims.(jwt.MapClaims)[claim].(string)
}
