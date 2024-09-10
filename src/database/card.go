package database

import (
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
)

type CardType string

const (
	JudgeCard  CardType = "Judge"
	PlayerCard CardType = "Player"
)

type Card struct {
	Id                uuid.UUID
	CreatedOnDate     time.Time
	ChangedOnDate     time.Time
	CreatedByPlayerId uuid.UUID
	ChangedByPlayerId uuid.UUID

	DeckId uuid.UUID
	Type   CardType
	Text   string
}

func GetCardsInDeck(deckId uuid.UUID, textSearch string, cardTypeSearch string) ([]Card, error) {
	if textSearch == "" {
		textSearch = "%"
	}

	if cardTypeSearch == "" {
		cardTypeSearch = "%"
	}

	sqlString := `
		SELECT
			ID,
			CREATED_ON_DATE,
			CHANGED_ON_DATE,
			CREATED_BY_PLAYER_ID,
			CHANGED_BY_PLAYER_ID,
			DECK_ID,
			TYPE,
			TEXT
		FROM CARD
		WHERE DECK_ID = ?
			AND TEXT LIKE ?
			AND TYPE LIKE ?
		ORDER BY TYPE ASC, CHANGED_ON_DATE DESC, TEXT ASC
	`
	rows, err := Query(sqlString, deckId, textSearch, cardTypeSearch)
	if err != nil {
		return nil, err
	}

	result := make([]Card, 0)
	for rows.Next() {
		var card Card
		if err := rows.Scan(
			&card.Id,
			&card.CreatedOnDate,
			&card.ChangedOnDate,
			&card.CreatedByPlayerId,
			&card.ChangedByPlayerId,
			&card.DeckId,
			&card.Type,
			&card.Text); err != nil {
			continue
		}
		result = append(result, card)
	}
	return result, nil
}

func GetCard(id uuid.UUID) (Card, error) {
	var card Card

	sqlString := `
		SELECT
			ID,
			CREATED_ON_DATE,
			CHANGED_ON_DATE,
			CREATED_BY_PLAYER_ID,
			CHANGED_BY_PLAYER_ID,
			DECK_ID,
			TYPE,
			TEXT
		FROM CARD
		WHERE ID = ?
	`
	rows, err := Query(sqlString, id)
	if err != nil {
		return card, err
	}

	for rows.Next() {
		if err := rows.Scan(
			&card.Id,
			&card.CreatedOnDate,
			&card.ChangedOnDate,
			&card.CreatedByPlayerId,
			&card.ChangedByPlayerId,
			&card.DeckId,
			&card.Type,
			&card.Text); err != nil {
			log.Println(err)
			return card, errors.New("failed to scan row in query results")
		}
	}

	return card, nil
}

func CreateCard(playerId uuid.UUID, deckId uuid.UUID, cardType CardType, text string) (uuid.UUID, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		log.Println(err)
		return id, errors.New("failed to generate new id")
	}

	sqlString := `
		INSERT INTO CARD (ID, CREATED_BY_PLAYER_ID, CHANGED_BY_PLAYER_ID, DECK_ID, TYPE, TEXT)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	return id, Execute(sqlString, id, playerId, playerId, deckId, cardType, text)
}

func GetCardId(deckId uuid.UUID, text string) (uuid.UUID, error) {
	var id uuid.UUID

	sqlString := `
		SELECT
			ID
		FROM CARD
		WHERE DECK_ID = ?
			AND TEXT = ?
	`
	rows, err := Query(sqlString, deckId, text)
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

func SetCardType(playerId uuid.UUID, id uuid.UUID, cardType CardType) error {
	sqlString := `
		UPDATE CARD
		SET
			TYPE = ?,
			CHANGED_BY_PLAYER_ID = ?
		WHERE ID = ?
	`
	return Execute(sqlString, cardType, playerId, id)
}

func SetCardText(playerId uuid.UUID, id uuid.UUID, text string) error {
	sqlString := `
		UPDATE CARD
		SET
			TEXT = ?,
			CHANGED_BY_PLAYER_ID = ?
		WHERE ID = ?
	`
	return Execute(sqlString, text, playerId, id)
}

func DeleteCard(id uuid.UUID) error {
	sqlString := `
		DELETE FROM CARD
		WHERE ID = ?
	`
	return Execute(sqlString, id)
}
