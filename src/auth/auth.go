package auth

import (
	"errors"
	"net/http"
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

const cookieNamePlayerToken string = "CARD-JUDGE-PLAYER-TOKEN"
const cookieNameRedirectURL string = "CARD-JUDGE-REDIRECT-URL"
const claimKey string = "playerName"

var jwtSecret []byte = []byte(os.Getenv("GFB_JWT_SECRET"))

func GetRedirectURL(r *http.Request) string {
	redirectPath := "/"
	for _, c := range r.Cookies() {
		if c.Name != cookieNameRedirectURL {
			continue
		}
		redirectPath = c.Value
		break
	}
	return redirectPath
}

func SetRedirectURL(w http.ResponseWriter, url string) {
	cookie := http.Cookie{
		Name:    cookieNameRedirectURL,
		Value:   url,
		Path:    "/",
		Expires: time.Now().Add(time.Hour * 12),
	}
	http.SetCookie(w, &cookie)
}

func RemoveRedirectURL(w http.ResponseWriter) {
	cookie := http.Cookie{
		Name:    cookieNameRedirectURL,
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),
	}
	http.SetCookie(w, &cookie)
}

func GetPlayerName(r *http.Request) (string, error) {
	cookieValue := ""
	for _, c := range r.Cookies() {
		if c.Name != cookieNamePlayerToken {
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

	cookie := http.Cookie{
		Name:    cookieNamePlayerToken,
		Value:   tokenString,
		Path:    "/",
		Expires: time.Now().Add(time.Hour * 12),
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
