package auth

import (
	"os"

	"github.com/borosr/realworld/lib/broken"
	"github.com/golang-jwt/jwt"
)

var (
	hmacSampleSecret       = os.Getenv("JWT_SIGNING_KEY")
	ErrUnableToVerifyToken = broken.Forbidden("unable to verify token")
)

func Verify(token string) (map[string]interface{}, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, broken.Forbiddenf("unexpected signing method: %v", token.Header["alg"])
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
	signedString, err := token.SignedString([]byte(hmacSampleSecret))
	if err != nil {
		return "", broken.Forbidden(err.Error())
	}
	return signedString, nil
}
