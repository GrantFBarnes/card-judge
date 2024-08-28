package auth

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const cookieNamePlayerToken string = "CARD-JUDGE-PLAYER-TOKEN"
const cookieNameAccessToken string = "CARD-JUDGE-ACCESS-TOKEN"
const claimKey string = "value"

var jwtSecret []byte = []byte(os.Getenv("GFB_JWT_SECRET"))

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

func RemovePlayerName(w http.ResponseWriter) {
	cookie := getRemovalCookie(cookieNamePlayerToken)
	http.SetCookie(w, &cookie)
}

func HasAccess(r *http.Request, id uuid.UUID) bool {
	accessIds, err := getAccessIds(r)
	if err != nil {
		return false
	}

	for _, v := range accessIds {
		if v == id {
			return true
		}
	}

	return false
}

func AddAccessId(w http.ResponseWriter, r *http.Request, accessId uuid.UUID) error {
	accessIds, err := getAccessIds(r)
	if err != nil {
		return err
	}

	alreadyAdded := false
	for _, v := range accessIds {
		if v == accessId {
			alreadyAdded = true
			break
		}
	}

	if !alreadyAdded {
		accessIds = append(accessIds, accessId)
	}

	return SetAccessIds(w, accessIds)
}

func getAccessIds(r *http.Request) ([]uuid.UUID, error) {
	cookieValue := ""
	for _, c := range r.Cookies() {
		if c.Name != cookieNameAccessToken {
			continue
		}
		cookieValue = c.Value
		break
	}

	var accessIds = make([]uuid.UUID, 0)
	if cookieValue == "" {
		return accessIds, nil
	}

	claimValue, err := getTokenClaimValue(cookieValue)
	if err != nil {
		return nil, err
	}

	accessStrings := strings.Split(claimValue, " ")

	accessIds = accessToUuid(accessStrings)

	return accessIds, nil
}

func SetAccessIds(w http.ResponseWriter, accessIds []uuid.UUID) error {
	accessStrings := accessToString(accessIds)
	tokenString, err := getTokenString(strings.Join(accessStrings, " "))
	if err != nil {
		return err
	}

	cookie := http.Cookie{
		Name:    cookieNameAccessToken,
		Value:   tokenString,
		Path:    "/",
		Expires: time.Now().Add(time.Hour * 12),
	}
	http.SetCookie(w, &cookie)
	return nil
}

func RemoveAccess(w http.ResponseWriter) {
	cookie := getRemovalCookie(cookieNameAccessToken)
	http.SetCookie(w, &cookie)
}

func accessToUuid(accessStrings []string) []uuid.UUID {
	accessIds := make([]uuid.UUID, len(accessStrings))
	for i := range accessStrings {
		id, err := uuid.Parse(accessStrings[i])
		if err != nil {
			continue
		}
		accessIds[i] = id
	}
	return accessIds
}

func accessToString(accessIds []uuid.UUID) []string {
	accessStrings := make([]string, len(accessIds))
	for i := range accessIds {
		accessStrings[i] = accessIds[i].String()
	}
	return accessStrings
}

func getRemovalCookie(cookieName string) http.Cookie {
	return http.Cookie{
		Name:    cookieName,
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),
	}
}

func getTokenString(tokenStringValue string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		claimKey: tokenStringValue,
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
