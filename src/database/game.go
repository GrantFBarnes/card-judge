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
			L.JUDGE_USER_ID,
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
			&data.JudgeUserId,
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

	lobbyUserId, err := getLobbyUserId(lobbyId, userId)
	if err != nil {
		return data, err
	}

	data.Cards, err = getLobbyUserHand(lobbyUserId)
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
			LU.USER_ID,
			U.NAME AS USER_NAME,
			COUNT(LR.ID) AS WINS
		FROM LOBBY_USER AS LU
			LEFT JOIN LOBBY_RESULT AS LR ON LR.LOBBY_USER_ID = LU.ID
			INNER JOIN USER AS U ON U.ID = LU.USER_ID
		WHERE LU.LOBBY_ID = ?
		GROUP BY LU.USER_ID
		ORDER BY COUNT(LR.ID) DESC, U.NAME ASC
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

func DrawLobbyUserHand(lobbyId uuid.UUID, userId uuid.UUID) (data lobbyGameUser, err error) {
	lobbyUserId, err := getLobbyUserId(lobbyId, userId)
	if err != nil {
		return data, err
	}

	handCount, err := getLobbyUserHandCount(lobbyUserId)
	if err != nil {
		return data, err
	}

	cardsToDraw := 8 - handCount
	if cardsToDraw > 0 {
		sqlString := `
			INSERT INTO LOBBY_USER_CARD
				(
					LOBBY_USER_ID,
					LOBBY_ID,
					USER_ID,
					CARD_ID
				)
			SELECT DISTINCT
				? AS LOBBY_USER_ID,
				? AS LOBBY_ID,
				? AS USER_ID,
				C.ID AS CARD_ID
			FROM DRAW_PILE AS DP
				INNER JOIN CARD AS C ON C.ID = DP.CARD_ID
				INNER JOIN CARD_TYPE AS CT ON CT.ID = C.CARD_TYPE_ID
			WHERE CT.NAME = 'Player'
				AND DP.LOBBY_ID = ?
			ORDER BY RAND()
			LIMIT ?
		`
		err = Execute(sqlString, lobbyUserId, lobbyId, userId, lobbyId, cardsToDraw)
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

func DiscardLobbyUserHand(lobbyId uuid.UUID, userId uuid.UUID) (data lobbyGameUser, err error) {
	lobbyUserId, err := getLobbyUserId(lobbyId, userId)
	if err != nil {
		return data, err
	}

	sqlString := `
		DELETE FROM LOBBY_USER_CARD
		WHERE LOBBY_USER_ID = ?
			AND LOBBY_ID = ?
			AND USER_ID = ?
	`
	err = Execute(sqlString, lobbyUserId, lobbyId, userId)
	if err != nil {
		return data, err
	}

	return GetLobbyGameUser(lobbyId, userId)
}

func DiscardLobbyUserCard(lobbyId uuid.UUID, userId uuid.UUID, cardId uuid.UUID) (data lobbyGameUser, err error) {
	lobbyUserId, err := getLobbyUserId(lobbyId, userId)
	if err != nil {
		return data, err
	}

	sqlString := `
		DELETE FROM LOBBY_USER_CARD
		WHERE LOBBY_USER_ID = ?
			AND LOBBY_ID = ?
			AND USER_ID = ?
			AND CARD_ID = ?
	`
	err = Execute(sqlString, lobbyUserId, lobbyId, userId, cardId)
	if err != nil {
		return data, err
	}

	return GetLobbyGameUser(lobbyId, userId)
}

func getLobbyUserId(lobbyId uuid.UUID, userId uuid.UUID) (lobbyUserId uuid.UUID, err error) {
	sqlString := `
		SELECT
			ID
		FROM LOBBY_USER
		WHERE LOBBY_ID = ?
			AND USER_ID = ?
		LIMIT 1
	`
	rows, err := Query(sqlString, lobbyId, userId)
	if err != nil {
		return lobbyUserId, err
	}

	for rows.Next() {
		if err := rows.Scan(&lobbyUserId); err != nil {
			return lobbyUserId, err
		}
	}

	return lobbyUserId, nil
}

func getLobbyUserHandCount(lobbyUserId uuid.UUID) (handCount int, err error) {
	sqlString := `
		SELECT
			COUNT(CARD_ID)
		FROM LOBBY_USER_CARD
		WHERE LOBBY_USER_ID = ?
	`
	rows, err := Query(sqlString, lobbyUserId)
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
			INNER JOIN LOBBY_USER_CARD AS LUC ON LUC.LOBBY_ID = DP.LOBBY_ID AND LUC.CARD_ID = DP.CARD_ID
	`
	return Execute(sqlString)
}

func getLobbyUserHand(lobbyUserId uuid.UUID) ([]Card, error) {
	sqlString := `
		SELECT
			C.ID,
			C.TEXT
		FROM LOBBY_USER_CARD AS LUC
			INNER JOIN CARD AS C ON C.ID = LUC.CARD_ID
		WHERE LUC.LOBBY_USER_ID = ?
		ORDER BY C.TEXT
	`
	rows, err := Query(sqlString, lobbyUserId)
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
