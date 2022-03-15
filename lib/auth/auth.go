package auth

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt"
)

var (
	//hmacSampleSecret       = os.Getenv("JWT_SIGNING_KEY")
	hmacSampleSecret       = "test_key_123"
	ErrUnableToVerifyToken = errors.New("unable to verify token")
)

func Verify(token string) (map[string]interface{}, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(hmacSampleSecret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return nil, ErrUnableToVerifyToken
	}

	return claims, nil
}

func Sign(claims map[string]interface{}) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(claims))
	return token.SignedString([]byte(hmacSampleSecret))
}
