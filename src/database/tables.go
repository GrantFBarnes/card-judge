package database

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Player struct {
	Id           uuid.UUID
	DateAdded    time.Time
	DateModified time.Time

	Name string
}

type Lobby struct {
	Id           uuid.UUID
	DateAdded    time.Time
	DateModified time.Time

	Name     string
	Password sql.NullString
}

type Deck struct {
	Id           uuid.UUID
	DateAdded    time.Time
	DateModified time.Time

	Name     string
	Password sql.NullString
}

type Card struct {
	Id           uuid.UUID
	DateAdded    time.Time
	DateModified time.Time

	DeckId uuid.UUID
	Type   string
	Text   string
}
