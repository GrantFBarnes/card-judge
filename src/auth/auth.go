package auth

import (
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/grantfbarnes/card-judge/helper"
)

func GetCookiePlayerId(r *http.Request) (uuid.UUID, error) {
	cookieValue, err := getCookie(r, cookieNamePlayerToken)
	if err != nil {
		return uuid.Nil, err
	}

	tokenValue, err := getTokenStringValue(cookieValue)
	if err != nil {
		return uuid.Nil, err
	}

	playerId, err := uuid.Parse(tokenValue)
	if err != nil {
		return uuid.Nil, err
	}

	return playerId, nil
}

func SetCookiePlayerId(w http.ResponseWriter, playerId uuid.UUID) error {
	tokenString, err := getValueTokenString(playerId.String())
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

func RemoveCookiePlayerId(w http.ResponseWriter) {
	cookie := getRemovalCookie(cookieNamePlayerToken)
	http.SetCookie(w, &cookie)
}

func HasCookieAccess(r *http.Request, id uuid.UUID) bool {
	accessIds, err := getCookieAccessIds(r)
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

func AddCookieAccessId(w http.ResponseWriter, r *http.Request, accessId uuid.UUID) error {
	accessIds, err := getCookieAccessIds(r)
	if err != nil {
		return err
	}

	for _, v := range accessIds {
		if v == accessId {
			// already have access
			return nil
		}
	}

	accessIds = append(accessIds, accessId)

	accessStrings := helper.ConvertArrayUuidsToStrings(accessIds)
	tokenString, err := getValueTokenString(strings.Join(accessStrings, " "))
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

func getCookieAccessIds(r *http.Request) ([]uuid.UUID, error) {
	var accessIds = make([]uuid.UUID, 0)

	cookieValue, err := getCookie(r, cookieNameAccessToken)
	if err != nil {
		return accessIds, nil
	}

	tokenValue, err := getTokenStringValue(cookieValue)
	if err != nil {
		return nil, err
	}

	accessStrings := strings.Split(tokenValue, " ")

	accessIds = helper.ConvertArrayStringsToUuids(accessStrings)

	return accessIds, nil
}
