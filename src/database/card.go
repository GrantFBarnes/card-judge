package database

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Card struct {
	Id           uuid.UUID
	DateAdded    time.Time
	DateModified time.Time

	DeckId uuid.UUID
	Type   string
	Text   string
}

func GetCards(dbcs string, deckId uuid.UUID) ([]Card, error) {
	db, err := sql.Open("mysql", dbcs)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	selectStatment, err := db.Prepare(`
		SELECT ID
			 , DATE_ADDED
			 , DATE_MODIFIED
			 , TYPE
			 , TEXT
	 	FROM CARD 
		WHERE DECK_ID = ?
	`)
	if err != nil {
		return nil, err
	}
	defer selectStatment.Close()

	rows, err := selectStatment.Query(deckId)
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
