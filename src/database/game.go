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
			&data.CardCount); err != nil {
			return data, err
		}
	}

	return data, nil
}

type lobbyGameUser struct {
	Lobby        Lobby
	User         User
	Cards        []Card
	MaxCardCount int
	CardCount    int
}

func GetLobbyGameUser(lobbyId uuid.UUID, userId uuid.UUID) (data lobbyGameUser, err error) {
	data.Lobby, err = GetLobby(lobbyId)
	if err != nil {
		return data, err
	}

	data.User, err = GetUser(userId)
	if err != nil {
		return data, err
	}

	playerId, err := getPlayerId(lobbyId, userId)
	if err != nil {
		return data, err
	}

	data.Cards, err = getPlayerHand(playerId)
	if err != nil {
		return data, err
	}

	data.MaxCardCount = 8
	data.CardCount = len(data.Cards)

	return data, nil
}

type lobbyGameStats struct {
	UserId   uuid.UUID
	UserName string
	Wins     int
}

func GetLobbyGameStats(lobbyId uuid.UUID) ([]lobbyGameStats, error) {
	sqlString := `
		SELECT
			P.USER_ID,
			U.NAME AS USER_NAME,
			COUNT(W.ID) AS WINS
		FROM PLAYER AS P
			LEFT JOIN WIN AS W ON W.PLAYER_ID = P.ID
			INNER JOIN USER AS U ON U.ID = P.USER_ID
		WHERE P.LOBBY_ID = ?
		GROUP BY P.USER_ID
		ORDER BY COUNT(W.ID) DESC, U.NAME ASC
	`
	rows, err := Query(sqlString, lobbyId)
	if err != nil {
		return nil, err
	}

	result := make([]lobbyGameStats, 0)
	for rows.Next() {
		var stats lobbyGameStats
		if err := rows.Scan(
			&stats.UserId,
			&stats.UserName,
			&stats.Wins); err != nil {
			continue
		}
		result = append(result, stats)
	}
	return result, nil
}

func DrawPlayerHand(lobbyId uuid.UUID, userId uuid.UUID) (data lobbyGameUser, err error) {
	playerId, err := getPlayerId(lobbyId, userId)
	if err != nil {
		return data, err
	}

	handCount, err := getPlayerHandCount(playerId)
	if err != nil {
		return data, err
	}

	cardsToDraw := 8 - handCount
	if cardsToDraw > 0 {
		sqlString := `
			INSERT INTO HAND
				(
					PLAYER_ID,
					CARD_ID
				)
			SELECT DISTINCT
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
		err = Execute(sqlString, playerId, lobbyId, cardsToDraw)
		if err != nil {
			return data, err
		}

		err = removeUserHandFromLobbyCards()
		if err != nil {
			return data, err
		}
	}

	return GetLobbyGameUser(lobbyId, userId)
}

func DiscardPlayerHand(lobbyId uuid.UUID, userId uuid.UUID) (data lobbyGameUser, err error) {
	playerId, err := getPlayerId(lobbyId, userId)
	if err != nil {
		return data, err
	}

	sqlString := `
		DELETE FROM HAND
		WHERE PLAYER_ID = ?
	`
	err = Execute(sqlString, playerId)
	if err != nil {
		return data, err
	}

	return GetLobbyGameUser(lobbyId, userId)
}

func DiscardPlayerCard(lobbyId uuid.UUID, userId uuid.UUID, cardId uuid.UUID) (data lobbyGameUser, err error) {
	playerId, err := getPlayerId(lobbyId, userId)
	if err != nil {
		return data, err
	}

	sqlString := `
		DELETE FROM HAND
		WHERE PLAYER_ID = ?
			AND CARD_ID = ?
	`
	err = Execute(sqlString, playerId, cardId)
	if err != nil {
		return data, err
	}

	return GetLobbyGameUser(lobbyId, userId)
}

func getPlayerId(lobbyId uuid.UUID, userId uuid.UUID) (playerId uuid.UUID, err error) {
	sqlString := `
		SELECT
			ID
		FROM PLAYER
		WHERE LOBBY_ID = ?
			AND USER_ID = ?
		LIMIT 1
	`
	rows, err := Query(sqlString, lobbyId, userId)
	if err != nil {
		return playerId, err
	}

	for rows.Next() {
		if err := rows.Scan(&playerId); err != nil {
			return playerId, err
		}
	}

	return playerId, nil
}

func getPlayerHandCount(playerId uuid.UUID) (handCount int, err error) {
	sqlString := `
		SELECT
			COUNT(CARD_ID)
		FROM HAND
		WHERE PLAYER_ID = ?
	`
	rows, err := Query(sqlString, playerId)
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

func removeUserHandFromLobbyCards() error {
	sqlString := `
		DELETE DP
		FROM DRAW_PILE AS DP
			INNER JOIN PLAYER AS P ON P.LOBBY_ID = DP.LOBBY_ID
			INNER JOIN HAND AS H ON H.PLAYER_ID = P.ID AND H.CARD_ID = DP.CARD_ID
	`
	return Execute(sqlString)
}

func getPlayerHand(playerId uuid.UUID) ([]Card, error) {
	sqlString := `
		SELECT
			C.ID,
			C.TEXT
		FROM HAND AS H
			INNER JOIN CARD AS C ON C.ID = H.CARD_ID
		WHERE H.PLAYER_ID = ?
		ORDER BY C.TEXT
	`
	rows, err := Query(sqlString, playerId)
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
