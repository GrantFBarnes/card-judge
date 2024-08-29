package auth

import (
	"errors"
	"os"

	jwt "github.com/golang-jwt/jwt/v5"
)

const claimKey string = "value"

var jwtSecret []byte = []byte(os.Getenv("GFB_JWT_SECRET"))

func getValueTokenString(value string) (tokenString string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		claimKey: value,
	})
	tokenString, err = token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func getTokenStringValue(tokenString string) (value string, err error) {
	token, err := getTokenStringToken(tokenString)
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims[claimKey].(string), nil
	} else {
		return "", errors.New("could not get token claims")
	}
}

func getTokenStringToken(tokenString string) (token *jwt.Token, err error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
}
