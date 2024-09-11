package database

import "github.com/google/uuid"

type LobbyGameInfo struct {
	Lobby
	CardCount int
}

func GetLobbyGameInfo(lobbyId uuid.UUID) (data LobbyGameInfo, err error) {
	sqlString := `
		SELECT
			L.ID,
			L.NAME,
			L.JUDGE_PLAYER_ID,
			COUNT(DP.CARD_ID) AS CARD_COUNT
		FROM LOBBY AS L
			INNER JOIN DRAW_PILE AS DP ON DP.LOBBY_ID = L.ID
		WHERE L.ID = ?
		GROUP BY L.ID
	`
	rows, err := Query(sqlString, lobbyId)
	if err != nil {
		return data, err
	}

	for rows.Next() {
		if err := rows.Scan(
			&data.Id,
			&data.Name,
			&data.JudgePlayerId,
			&data.CardCount); err != nil {
			return data, err
		}
	}

	return data, nil
}

type lobbyGamePlayer struct {
	Lobby        Lobby
	Player       Player
	Cards        []Card
	MaxCardCount int
	CardCount    int
}

func GetLobbyGamePlayer(lobbyId uuid.UUID, playerId uuid.UUID) (data lobbyGamePlayer, err error) {
	data.Lobby, err = GetLobby(lobbyId)
	if err != nil {
		return data, err
	}

	data.Player, err = GetPlayer(playerId)
	if err != nil {
		return data, err
	}

	lobbyPlayerId, err := getLobbyPlayerId(lobbyId, playerId)
	if err != nil {
		return data, err
	}

	data.Cards, err = getLobbyPlayerHand(lobbyPlayerId)
	if err != nil {
		return data, err
	}

	data.MaxCardCount = 8
	data.CardCount = len(data.Cards)

	return data, nil
}

type lobbyGameStats struct {
	PlayerId   uuid.UUID
	PlayerName string
	Wins       int
}

func GetLobbyGameStats(lobbyId uuid.UUID) ([]lobbyGameStats, error) {
	sqlString := `
		SELECT
			LP.PLAYER_ID,
			P.NAME AS PLAYER_NAME,
			COUNT(LR.ID) AS WINS
		FROM LOBBY_PLAYER AS LP
			LEFT JOIN LOBBY_RESULT AS LR ON LR.LOBBY_PLAYER_ID = LP.ID
			INNER JOIN PLAYER AS P ON P.ID = LP.PLAYER_ID
		WHERE LP.LOBBY_ID = ?
		GROUP BY LP.PLAYER_ID
		ORDER BY COUNT(LR.ID) DESC, P.NAME ASC
	`
	rows, err := Query(sqlString, lobbyId)
	if err != nil {
		return nil, err
	}

	result := make([]lobbyGameStats, 0)
	for rows.Next() {
		var stats lobbyGameStats
		if err := rows.Scan(
			&stats.PlayerId,
			&stats.PlayerName,
			&stats.Wins); err != nil {
			continue
		}
		result = append(result, stats)
	}
	return result, nil
}

func DrawLobbyPlayerHand(lobbyId uuid.UUID, playerId uuid.UUID) (data lobbyGamePlayer, err error) {
	lobbyPlayerId, err := getLobbyPlayerId(lobbyId, playerId)
	if err != nil {
		return data, err
	}

	handCount, err := getLobbyPlayerHandCount(lobbyPlayerId)
	if err != nil {
		return data, err
	}

	cardsToDraw := 8 - handCount
	if cardsToDraw > 0 {
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
			FROM DRAW_PILE AS DP
				INNER JOIN CARD AS C ON C.ID = DP.CARD_ID
				INNER JOIN CARD_TYPE AS CT ON CT.ID = C.CARD_TYPE_ID
			WHERE CT.NAME = 'Player'
				AND DP.LOBBY_ID = ?
			ORDER BY RAND()
			LIMIT ?
		`
		err = Execute(sqlString, lobbyPlayerId, lobbyId, playerId, lobbyId, cardsToDraw)
		if err != nil {
			return data, err
		}

		err = removePlayerHandFromLobbyCards()
		if err != nil {
			return data, err
		}
	}

	return GetLobbyGamePlayer(lobbyId, playerId)
}

func DiscardLobbyPlayerHand(lobbyId uuid.UUID, playerId uuid.UUID) (data lobbyGamePlayer, err error) {
	lobbyPlayerId, err := getLobbyPlayerId(lobbyId, playerId)
	if err != nil {
		return data, err
	}

	sqlString := `
		DELETE FROM LOBBY_PLAYER_CARD
		WHERE LOBBY_PLAYER_ID = ?
			AND LOBBY_ID = ?
			AND PLAYER_ID = ?
	`
	err = Execute(sqlString, lobbyPlayerId, lobbyId, playerId)
	if err != nil {
		return data, err
	}

	return GetLobbyGamePlayer(lobbyId, playerId)
}

func DiscardLobbyPlayerCard(lobbyId uuid.UUID, playerId uuid.UUID, cardId uuid.UUID) (data lobbyGamePlayer, err error) {
	lobbyPlayerId, err := getLobbyPlayerId(lobbyId, playerId)
	if err != nil {
		return data, err
	}

	sqlString := `
		DELETE FROM LOBBY_PLAYER_CARD
		WHERE LOBBY_PLAYER_ID = ?
			AND LOBBY_ID = ?
			AND PLAYER_ID = ?
			AND CARD_ID = ?
	`
	err = Execute(sqlString, lobbyPlayerId, lobbyId, playerId, cardId)
	if err != nil {
		return data, err
	}

	return GetLobbyGamePlayer(lobbyId, playerId)
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
		DELETE DP
		FROM DRAW_PILE AS DP
			INNER JOIN LOBBY_PLAYER_CARD AS LPC ON LPC.LOBBY_ID = DP.LOBBY_ID AND LPC.CARD_ID = DP.CARD_ID
	`
	return Execute(sqlString)
}

func getLobbyPlayerHand(lobbyPlayerId uuid.UUID) ([]Card, error) {
	sqlString := `
		SELECT
			C.ID,
			C.TEXT
		FROM LOBBY_PLAYER_CARD AS LPC
			INNER JOIN CARD AS C ON C.ID = LPC.CARD_ID
		WHERE LPC.LOBBY_PLAYER_ID = ?
		ORDER BY C.TEXT
	`
	rows, err := Query(sqlString, lobbyPlayerId)
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
