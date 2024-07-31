package database

import (
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
	Password string
}

type Deck struct {
	Id           uuid.UUID
	DateAdded    time.Time
	DateModified time.Time

	Name     string
	Password string
}

type Card struct {
	Id           uuid.UUID
	DateAdded    time.Time
	DateModified time.Time

	DeckId uuid.UUID
	Type   string
	Text   string
}
