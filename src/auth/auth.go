package auth

import (
	"errors"
	"os"

	jwt "github.com/golang-jwt/jwt/v5"
)

var secret []byte = []byte(os.Getenv("GFB_JWT_SECRET"))

func GetTokenString(playerName string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"playerName": playerName,
	})
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func GetTokenClaims(tokenString string) (map[string]string, error) {
	verified, err := getToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims, ok := verified.Claims.(jwt.MapClaims); ok {
		var result = make(map[string]string)
		result["playerName"] = claims["playerName"].(string)
		return result, nil
	} else {
		return nil, errors.New("could not get claims")
	}
}

func getToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
}
