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

const claimKey string = "value"

var jwtSecret []byte = []byte(os.Getenv("GFB_JWT_SECRET"))

func GetPlayerId(r *http.Request) (uuid.UUID, error) {
	cookieValue, err := getCookie(r, cookieNamePlayerToken)
	if err != nil {
		return uuid.Nil, err
	}

	claimValue, err := getTokenClaimValue(cookieValue)
	if err != nil {
		return uuid.Nil, err
	}

	result, err := uuid.Parse(claimValue)
	if err != nil {
		return uuid.Nil, err
	}

	return result, nil
}

func SetPlayerId(w http.ResponseWriter, playerId uuid.UUID) error {
	tokenString, err := getTokenString(playerId.String())
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

func RemovePlayerId(w http.ResponseWriter) {
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
	var accessIds = make([]uuid.UUID, 0)

	cookieValue, err := getCookie(r, cookieNameAccessToken)
	if err != nil {
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
