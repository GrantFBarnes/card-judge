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
	Id                uuid.UUID
	CreatedOnDate     time.Time
	ChangedOnDate     time.Time
	CreatedByPlayerId uuid.UUID
	ChangedByPlayerId uuid.UUID

	Name         string
	PasswordHash sql.NullString
	InLobby      bool
}

func GetDecks() ([]Deck, error) {
	rows, err := Query(`
		SELECT
			ID,
			CREATED_ON_DATE,
			CHANGED_ON_DATE,
			CREATED_BY_PLAYER_ID,
			CHANGED_BY_PLAYER_ID,
			NAME,
			PASSWORD_HASH
		FROM DECK
		ORDER BY CHANGED_ON_DATE DESC
	`)
	if err != nil {
		return nil, err
	}

	result := make([]Deck, 0)
	for rows.Next() {
		var deck Deck
		if err := rows.Scan(
			&deck.Id,
			&deck.CreatedOnDate,
			&deck.ChangedOnDate,
			&deck.CreatedByPlayerId,
			&deck.ChangedByPlayerId,
			&deck.Name,
			&deck.PasswordHash); err != nil {
			continue
		}
		result = append(result, deck)
	}
	return result, nil
}

func SearchDecks(search string) ([]Deck, error) {
	rows, err := Query(`
		SELECT
			ID,
			CREATED_ON_DATE,
			CHANGED_ON_DATE,
			CREATED_BY_PLAYER_ID,
			CHANGED_BY_PLAYER_ID,
			NAME,
			PASSWORD_HASH
		FROM DECK
		WHERE NAME LIKE ?
		ORDER BY CHANGED_ON_DATE DESC
	`, search)
	if err != nil {
		return nil, err
	}

	result := make([]Deck, 0)
	for rows.Next() {
		var deck Deck
		if err := rows.Scan(
			&deck.Id,
			&deck.CreatedOnDate,
			&deck.ChangedOnDate,
			&deck.CreatedByPlayerId,
			&deck.ChangedByPlayerId,
			&deck.Name,
			&deck.PasswordHash); err != nil {
			continue
		}
		result = append(result, deck)
	}
	return result, nil
}

func GetDeck(id uuid.UUID) (Deck, error) {
	var deck Deck

	rows, err := Query(`
		SELECT
			ID,
			CREATED_ON_DATE,
			CHANGED_ON_DATE,
			CREATED_BY_PLAYER_ID,
			CHANGED_BY_PLAYER_ID,
			NAME,
			PASSWORD_HASH
		FROM DECK
		WHERE ID = ?
	`, id)
	if err != nil {
		return deck, err
	}

	for rows.Next() {
		if err := rows.Scan(
			&deck.Id,
			&deck.CreatedOnDate,
			&deck.ChangedOnDate,
			&deck.CreatedByPlayerId,
			&deck.ChangedByPlayerId,
			&deck.Name,
			&deck.PasswordHash); err != nil {
			log.Println(err)
			return deck, errors.New("failed to scan row in query results")
		}
	}

	return deck, nil
}

func CreateDeck(playerId uuid.UUID, name string, password string) (uuid.UUID, error) {
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
		INSERT INTO DECK (ID, CREATED_BY_PLAYER_ID, CHANGED_BY_PLAYER_ID, NAME, PASSWORD_HASH)
		VALUES (?, ?, ?, ?, ?)
	`
	if password == "" {
		return id, Execute(sqlString, id, playerId, playerId, name, nil)
	} else {
		return id, Execute(sqlString, id, playerId, playerId, name, passwordHash)
	}
}

func GetDeckId(name string) (uuid.UUID, error) {
	var id uuid.UUID

	rows, err := Query(`
		SELECT
			ID
		FROM DECK
		WHERE NAME = ?
	`, name)
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

func SetDeckName(playerId uuid.UUID, id uuid.UUID, name string) error {
	sqlString := `
		UPDATE DECK
		SET
			NAME = ?,
			CHANGED_ON_DATE = CURRENT_TIMESTAMP(),
			CHANGED_BY_PLAYER_ID = ?
		WHERE ID = ?
	`
	return Execute(sqlString, name, playerId, id)
}

func SetDeckPassword(playerId uuid.UUID, id uuid.UUID, password string) error {
	passwordHash, err := auth.GetPasswordHash(password)
	if err != nil {
		log.Println(err)
		return errors.New("failed to hash password")
	}

	sqlString := `
		UPDATE DECK
		SET
			PASSWORD_HASH = ?,
			CHANGED_ON_DATE = CURRENT_TIMESTAMP(),
			CHANGED_BY_PLAYER_ID = ?
		WHERE ID = ?
	`
	if password == "" {
		return Execute(sqlString, nil, playerId, id)
	} else {
		return Execute(sqlString, passwordHash, playerId, id)
	}
}

func DeleteDeck(id uuid.UUID) error {
	sqlString := `
		DELETE FROM DECK
		WHERE ID = ?
	`
	return Execute(sqlString, id)
}
