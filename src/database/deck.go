package database

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/grantfbarnes/card-judge/auth"
)

type Deck struct {
	Id            uuid.UUID
	CreatedOnDate time.Time
	ChangedOnDate time.Time

	Name             string
	PasswordHash     sql.NullString
	IsPublicReadOnly bool
}

type DeckDetails struct {
	Deck
	CardCount int
}

func SearchDecks(search string) ([]DeckDetails, error) {
	if search == "" {
		search = "%"
	}

	sqlString := `
		SELECT
			D.ID,
			D.CREATED_ON_DATE,
			D.CHANGED_ON_DATE,
			D.NAME,
			D.PASSWORD_HASH,
			COUNT(C.ID) AS CARD_COUNT,
			D.IS_PUBLIC_READONLY
		FROM DECK AS D
			LEFT JOIN CARD AS C ON C.DECK_ID = D.ID
		WHERE D.NAME LIKE ?
		AND NOT IS_HIDDEN
		GROUP BY D.ID
		ORDER BY
			TO_DAYS(D.CHANGED_ON_DATE) DESC,
			D.NAME ASC,
			COUNT(C.ID) DESC
	`
	rows, err := query(sqlString, search)
	if err != nil {
		return nil, err
	}

	result := make([]DeckDetails, 0)
	for rows.Next() {
		var deckDetails DeckDetails
		if err := rows.Scan(
			&deckDetails.Id,
			&deckDetails.CreatedOnDate,
			&deckDetails.ChangedOnDate,
			&deckDetails.Name,
			&deckDetails.PasswordHash,
			&deckDetails.CardCount,
			&deckDetails.IsPublicReadOnly); err != nil {
			log.Println(err)
			return result, errors.New("failed to scan row in query results")
		}
		result = append(result, deckDetails)
	}
	return result, nil
}

func GetUserDecks(userId uuid.UUID) ([]Deck, error) {
	sqlString := "CALL SP_GET_DECK_ACCESS (?, ?)"
	rows, err := query(sqlString, userId, "LOBBY")
	if err != nil {
		return nil, err
	}

	result := make([]Deck, 0)
	for rows.Next() {
		var deck Deck
		if err := rows.Scan(
			&deck.Id,
			&deck.Name); err != nil {
			log.Println(err)
			return result, errors.New("failed to scan row in query results")
		}
		result = append(result, deck)
	}
	return result, nil
}

func GetDeck(id uuid.UUID) (Deck, error) {
	var deck Deck

	sqlString := `
		SELECT
			ID,
			CREATED_ON_DATE,
			CHANGED_ON_DATE,
			NAME,
			PASSWORD_HASH,
			IS_PUBLIC_READONLY
		FROM DECK
		WHERE ID = ?
	`
	rows, err := query(sqlString, id)
	if err != nil {
		return deck, err
	}

	for rows.Next() {
		if err := rows.Scan(
			&deck.Id,
			&deck.CreatedOnDate,
			&deck.ChangedOnDate,
			&deck.Name,
			&deck.PasswordHash,
			&deck.IsPublicReadOnly); err != nil {
			log.Println(err)
			return deck, errors.New("failed to scan row in query results")
		}
	}

	return deck, nil
}

func GetDeckPasswordHash(id uuid.UUID) (sql.NullString, error) {
	var passwordHash sql.NullString

	sqlString := `
		SELECT
			PASSWORD_HASH
		FROM DECK
		WHERE ID = ?
	`
	rows, err := query(sqlString, id)
	if err != nil {
		return passwordHash, err
	}

	for rows.Next() {
		if err := rows.Scan(&passwordHash); err != nil {
			log.Println(err)
			return passwordHash, errors.New("failed to scan row in query results")
		}
	}

	return passwordHash, nil
}

func CreateDeck(name string, password string, isPublicReadOnly bool) (uuid.UUID, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		log.Println(err)
		return id, errors.New("failed to generate new id")
	}

	passwordHash, err := auth.GetPasswordHash(password)
	if err != nil {
		log.Println(err)
		return id, errors.New("failed to hash password")
	}

	sqlString := `
		INSERT INTO DECK (ID, NAME, PASSWORD_HASH, IS_PUBLIC_READONLY)
		VALUES (?, ?, ?, ?)
	`
	if password == "" {
		return id, execute(sqlString, id, name, nil, isPublicReadOnly)
	} else {
		return id, execute(sqlString, id, name, passwordHash, isPublicReadOnly)
	}
}

func GetDeckId(name string) (uuid.UUID, error) {
	var id uuid.UUID

	sqlString := `
		SELECT
			ID
		FROM DECK
		WHERE NAME = ?
	`
	rows, err := query(sqlString, name)
	if err != nil {
		return id, err
	}

	for rows.Next() {
		if err := rows.Scan(&id); err != nil {
			log.Println(err)
			return id, errors.New("failed to scan row in query results")
		}
	}

	return id, nil
}

func SetDeckName(id uuid.UUID, name string) error {
	sqlString := `
		UPDATE DECK
		SET
			NAME = ?
		WHERE ID = ?
	`
	return execute(sqlString, name, id)
}

func SetIsPublicReadOnly(id uuid.UUID, isPublicReadOnly bool) error {
	sqlString := `
		UPDATE DECK
		SET
			IS_PUBLIC_READONLY = ?
		WHERE ID = ?
	`
	return execute(sqlString, isPublicReadOnly, id)
}

func SetDeckPassword(id uuid.UUID, password string) error {
	passwordHash, err := auth.GetPasswordHash(password)
	if err != nil {
		log.Println(err)
		return errors.New("failed to hash password")
	}

	sqlString := `
		UPDATE DECK
		SET
			PASSWORD_HASH = ?
		WHERE ID = ?
	`
	if password == "" {
		return execute(sqlString, nil, id)
	} else {
		return execute(sqlString, passwordHash, id)
	}
}

func DeleteDeck(id uuid.UUID) error {
	sqlString := `
		DELETE FROM DECK
		WHERE ID = ?
	`
	return execute(sqlString, id)
}
