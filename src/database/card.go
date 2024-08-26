package database

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type CardType string

const (
	Judge  CardType = "Judge"
	Player CardType = "Player"
)

type Card struct {
	Id           uuid.UUID
	DateAdded    time.Time
	DateModified time.Time

	DeckId uuid.UUID
	Type   CardType
	Text   string
}

func GetCardsInDeck(dbcs string, deckId uuid.UUID) ([]Card, error) {
	db, err := sql.Open("mysql", dbcs)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	statment, err := db.Prepare(`
		SELECT ID
			 , DATE_ADDED
			 , DATE_MODIFIED
			 , TYPE
			 , TEXT
	 	FROM CARD
		WHERE DECK_ID = ?
		ORDER BY TYPE, DATE_MODIFIED DESC
	`)
	if err != nil {
		return nil, err
	}
	defer statment.Close()

	rows, err := statment.Query(deckId)
	if err != nil {
		return nil, err
	}

	result := make([]Card, 0)
	for rows.Next() {
		var card Card
		if err := rows.Scan(
			&card.Id,
			&card.DateAdded,
			&card.DateModified,
			&card.Type,
			&card.Text); err != nil {
			continue
		}
		result = append(result, card)
	}
	return result, nil
}

func GetCard(dbcs string, id uuid.UUID) (Card, error) {
	var card Card

	db, err := sql.Open("mysql", dbcs)
	if err != nil {
		return card, err
	}
	defer db.Close()

	statment, err := db.Prepare(`
		SELECT ID
			 , DATE_ADDED
			 , DATE_MODIFIED
			 , DECK_ID
			 , TYPE
			 , TEXT
	 	FROM CARD
		WHERE ID = ?
	`)
	if err != nil {
		return card, err
	}
	defer statment.Close()

	rows, err := statment.Query(id)
	if err != nil {
		return card, err
	}

	for rows.Next() {
		if err := rows.Scan(
			&card.Id,
			&card.DateAdded,
			&card.DateModified,
			&card.DeckId,
			&card.Type,
			&card.Text); err != nil {
			return card, err
		}
	}

	return card, nil
}

func CreateCard(dbcs string, deckId uuid.UUID, cardType CardType, text string) (uuid.UUID, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return id, err
	}

	db, err := sql.Open("mysql", dbcs)
	if err != nil {
		return id, err
	}
	defer db.Close()

	statment, err := db.Prepare(`
		INSERT INTO CARD (ID, DECK_ID, TYPE, TEXT)
		VALUES (?, ?, ?, ?)
	`)
	if err != nil {
		return id, err
	}
	defer statment.Close()

	_, err = statment.Exec(id, deckId, cardType, text)
	if err != nil {
		return id, err
	}

	return id, nil
}

func UpdateCard(dbcs string, id uuid.UUID, text string) error {
	db, err := sql.Open("mysql", dbcs)
	if err != nil {
		return err
	}
	defer db.Close()

	statment, err := db.Prepare(`
		UPDATE CARD
		SET TEXT = ?
		WHERE ID = ?
	`)
	if err != nil {
		return err
	}
	defer statment.Close()

	_, err = statment.Exec(text, id)
	if err != nil {
		return err
	}

	return nil
}

func DeleteCard(dbcs string, id uuid.UUID) error {
	db, err := sql.Open("mysql", dbcs)
	if err != nil {
		return err
	}
	defer db.Close()

	statment, err := db.Prepare(`
		DELETE FROM CARD
		WHERE ID = ?
	`)
	if err != nil {
		return err
	}
	defer statment.Close()

	_, err = statment.Exec(id)
	if err != nil {
		return err
	}

	return nil
}
