package database

import "github.com/google/uuid"

func DrawLobbyPlayerHand(lobbyId uuid.UUID, playerId uuid.UUID) error {
	lobbyPlayerId, err := getLobbyPlayerId(lobbyId, playerId)
	if err != nil {
		return err
	}

	handCount, err := getLobbyPlayerHandCount(lobbyPlayerId)
	if err != nil {
		return err
	}

	cardsToDraw := 8 - handCount
	if cardsToDraw <= 0 {
		return nil
	}

	sqlString := `
		INSERT INTO LOBBY_PLAYER_CARD
			(
				LOBBY_PLAYER_ID,
				LOBBY_ID,
				PLAYER_ID,
				CARD_ID
			)
		SELECT DISTINCT
			? AS LOBBY_PLAYER_ID,
			? AS LOBBY_ID,
			? AS PLAYER_ID,
			C.ID AS CARD_ID
		FROM LOBBY_CARD AS LC
			INNER JOIN CARD AS C ON C.ID = LC.CARD_ID
		WHERE C.TYPE = 'Player'
			AND LC.LOBBY_ID = ?
		ORDER BY RAND()
		LIMIT ?
	`
	err = Execute(sqlString, lobbyPlayerId, lobbyId, playerId, lobbyId, cardsToDraw)
	if err != nil {
		return err
	}

	err = removePlayerHandFromLobbyCards()
	if err != nil {
		return err
	}

	return nil
}

func getLobbyPlayerId(lobbyId uuid.UUID, playerId uuid.UUID) (lobbyPlayerId uuid.UUID, err error) {
	sqlString := `
		SELECT
			ID
		FROM LOBBY_PLAYER
		WHERE LOBBY_ID = ?
			AND PLAYER_ID = ?
		LIMIT 1
	`
	rows, err := Query(sqlString, lobbyId, playerId)
	if err != nil {
		return lobbyPlayerId, err
	}

	for rows.Next() {
		if err := rows.Scan(&lobbyPlayerId); err != nil {
			return lobbyPlayerId, err
		}
	}

	return lobbyPlayerId, nil
}

func getLobbyPlayerHandCount(lobbyPlayerId uuid.UUID) (handCount int, err error) {
	sqlString := `
		SELECT
			COUNT(CARD_ID)
		FROM LOBBY_PLAYER_CARD
		WHERE LOBBY_PLAYER_ID = ?
	`
	rows, err := Query(sqlString, lobbyPlayerId)
	if err != nil {
		return handCount, err
	}

	for rows.Next() {
		if err := rows.Scan(&handCount); err != nil {
			return handCount, err
		}
	}

	return handCount, nil
}

func removePlayerHandFromLobbyCards() error {
	sqlString := `
		DELETE LC
		FROM LOBBY_CARD AS LC
			INNER JOIN LOBBY_PLAYER_CARD AS LPC ON LPC.LOBBY_ID = LC.LOBBY_ID AND LPC.CARD_ID = LC.CARD_ID
	`
	return Execute(sqlString)
}

func GetLobbyPlayerHand(lobbyId uuid.UUID, playerId uuid.UUID) ([]Card, error) {
	sqlString := `
		SELECT
			C.ID,
			C.TEXT
		FROM LOBBY_PLAYER_CARD AS LPC
			INNER JOIN CARD AS C ON C.ID = LPC.CARD_ID
		WHERE LPC.LOBBY_ID = ?
			AND LPC.PLAYER_ID = ?
		ORDER BY C.TEXT
	`
	rows, err := Query(sqlString, lobbyId, playerId)
	if err != nil {
		return nil, err
	}

	result := make([]Card, 0)
	for rows.Next() {
		var card Card
		if err := rows.Scan(
			&card.Id,
			&card.Text); err != nil {
			continue
		}
		result = append(result, card)
	}
	return result, nil
}

func GetLobbyCardCount(lobbyId uuid.UUID) (count int, err error) {
	sqlString := `
		SELECT
			COUNT(CARD_ID)
		FROM LOBBY_CARD
		WHERE LOBBY_ID = ?
	`
	rows, err := Query(sqlString, lobbyId)
	if err != nil {
		return count, err
	}

	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			return count, err
		}
	}

	return count, nil
}
