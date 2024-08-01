package auth

import (
	"errors"
	"net/http"
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

const cookieName string = "CARD-JUDGE-JWT-TOKEN"
const claimKey string = "playerName"

var jwtSecret []byte = []byte(os.Getenv("GFB_JWT_SECRET"))

func GetPlayerName(r *http.Request) (string, error) {
	cookieValue := ""
	for _, c := range r.Cookies() {
		if c.Name != cookieName {
			continue
		}
		cookieValue = c.Value
		break
	}

	if cookieValue == "" {
		return "", errors.New("cookie not found")
	}

	claimValue, err := getTokenClaimValue(cookieValue)
	if err != nil {
		return "", err
	}

	return claimValue, nil
}

func SetPlayerName(w http.ResponseWriter, playerName string) error {
	tokenString, err := getTokenString(playerName)
	if err != nil {
		return err
	}

	expiration := time.Now().Add(time.Hour * 12)
	cookie := http.Cookie{
		Name:    cookieName,
		Value:   tokenString,
		Expires: expiration,
	}
	http.SetCookie(w, &cookie)
	return nil
}

func getTokenString(playerName string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		claimKey: playerName,
	})
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func getTokenClaimValue(tokenString string) (string, error) {
	verified, err := getToken(tokenString)
	if err != nil {
		return "", err
	}

	if claims, ok := verified.Claims.(jwt.MapClaims); ok {
		return claims[claimKey].(string), nil
	} else {
		return "", errors.New("could not get claims")
	}
}

func getToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
}
