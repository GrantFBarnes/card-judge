package auth

import (
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/grantfbarnes/card-judge/helper"
)

func GetPlayerId(r *http.Request) (uuid.UUID, error) {
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

func SetPlayerId(w http.ResponseWriter, playerId uuid.UUID) error {
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

	tokenValue, err := getTokenStringValue(cookieValue)
	if err != nil {
		return nil, err
	}

	accessStrings := strings.Split(tokenValue, " ")

	accessIds = helper.ConvertArrayStringsToUuids(accessStrings)

	return accessIds, nil
}

func SetAccessIds(w http.ResponseWriter, accessIds []uuid.UUID) error {
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

func RemoveAccess(w http.ResponseWriter) {
	cookie := getRemovalCookie(cookieNameAccessToken)
	http.SetCookie(w, &cookie)
}
