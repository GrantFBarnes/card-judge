package database

import (
	"errors"
	"log"
	"reflect"

	"github.com/google/uuid"
)

type StatPersonal struct {
	GamePlayCount    int
	GameWinCount     int
	RoundPlayCount   int
	RoundWinCount    int
	CardPlayCount    int
	CardDrawCount    int
	CardDiscardCount int
	CardSkipCount    int
	LobbyKickCount   int
}

func GetPersonalStats(userId uuid.UUID) (StatPersonal, error) {
	var result StatPersonal

	sqlString := `
		SELECT
			(
				SELECT COUNT(DISTINCT LOBBY_ID)
				FROM LOG_RESPONSE_CARD
				WHERE PLAYER_USER_ID = U.ID
			) AS GAME_PLAY_COUNT,
			(
				SELECT
					COUNT(*)
				FROM (
						SELECT
							LRC.LOBBY_ID,
							LRC.PLAYER_USER_ID,
							COUNT(LW.ID) AS ROUND_WIN_COUNT,
							RANK() OVER (PARTITION BY LRC.LOBBY_ID ORDER BY ROUND_WIN_COUNT DESC) AS RANKING
						FROM LOG_RESPONSE_CARD AS LRC
							INNER JOIN LOG_WIN AS LW ON LW.RESPONSE_ID = LRC.RESPONSE_ID
						GROUP BY LRC.LOBBY_ID, LRC.PLAYER_USER_ID
					) AS ROUND_WINS
				WHERE PLAYER_USER_ID = U.ID
					AND RANKING = 1
			) AS GAME_WIN_COUNT,
			(
				SELECT COUNT(DISTINCT ROUND_ID)
				FROM LOG_RESPONSE_CARD
				WHERE PLAYER_USER_ID = U.ID
			) AS ROUND_PLAY_COUNT,
			(
				SELECT COUNT(DISTINCT LW.ID)
				FROM LOG_RESPONSE_CARD AS LRC
						INNER JOIN LOG_WIN AS LW ON LW.RESPONSE_ID = LRC.RESPONSE_ID
				WHERE LRC.PLAYER_USER_ID = U.ID
			) AS ROUND_WIN_COUNT,
			(
				SELECT COUNT(*)
				FROM LOG_RESPONSE_CARD
				WHERE PLAYER_USER_ID = U.ID
			) AS CARD_PLAY_COUNT,
			(SELECT COUNT(*) FROM LOG_DRAW WHERE USER_ID = U.ID) AS CARD_DRAW_COUNT,
			(SELECT COUNT(*) FROM LOG_DISCARD WHERE USER_ID = U.ID) AS CARD_DISCARD_COUNT,
			(SELECT COUNT(*) FROM LOG_SKIP WHERE USER_ID = U.ID) AS CARD_SKIP_COUNT,
			(SELECT COUNT(*) FROM LOG_KICK WHERE USER_ID = U.ID) AS LOBBY_KICK_COUNT
		FROM USER AS U
		WHERE U.ID = ?
	`
	rows, err := query(sqlString, userId)
	if err != nil {
		return result, err
	}

	for rows.Next() {
		if err := rows.Scan(
			&result.GamePlayCount,
			&result.GameWinCount,
			&result.RoundPlayCount,
			&result.RoundWinCount,
			&result.CardPlayCount,
			&result.CardDrawCount,
			&result.CardDiscardCount,
			&result.CardSkipCount,
			&result.LobbyKickCount); err != nil {
			log.Println(err)
			return result, errors.New("failed to scan row in query results")
		}
	}

	return result, nil
}

func GetLeaderboardStats(userId uuid.UUID, topic string, subject string) ([]string, [][]string, error) {
	resultHeaders := make([]string, 0)
	resultRows := make([][]string, 0)
	params := make([]any, 0)

	var sqlString string
	switch topic {
	case "game-win-ratio":
		switch subject {
		case "player":
			resultHeaders = append(resultHeaders, "Games Played")
			resultHeaders = append(resultHeaders, "Games Won")
			resultHeaders = append(resultHeaders, "Win Ratio")
			resultHeaders = append(resultHeaders, "Player")
			sqlString = `
				SELECT
					GP.GAME_PLAY_COUNT AS PLAY_COUNT,
					COALESCE(GW.GAME_WIN_COUNT, 0) AS WIN_COUNT,
					COALESCE((GW.GAME_WIN_COUNT * 1.0) / (GP.GAME_PLAY_COUNT * 1.0), 0.0) AS WIN_RATIO,
					U.NAME AS NAME
				FROM USER AS U
						INNER JOIN (
							SELECT
								PLAYER_USER_ID,
								COUNT(DISTINCT LOBBY_ID) AS GAME_PLAY_COUNT
							FROM LOG_RESPONSE_CARD
							GROUP BY PLAYER_USER_ID
						) AS GP ON GP.PLAYER_USER_ID = U.ID
						LEFT JOIN (
							SELECT
								PLAYER_USER_ID,
								COUNT(*) AS GAME_WIN_COUNT
							FROM (
								SELECT
									LRC.PLAYER_USER_ID,
									COUNT(LRC.ID) AS ROUND_WIN_COUNT,
									RANK() OVER (PARTITION BY LRC.LOBBY_ID ORDER BY ROUND_WIN_COUNT DESC) AS RANKING
								FROM LOG_RESPONSE_CARD AS LRC
									INNER JOIN LOG_WIN AS LW ON LW.RESPONSE_ID = LRC.RESPONSE_ID
								GROUP BY LRC.LOBBY_ID, LRC.PLAYER_USER_ID
							) AS GAME_RANK
							WHERE RANKING = 1
							GROUP BY PLAYER_USER_ID
						) AS GW ON GW.PLAYER_USER_ID = U.ID
				ORDER BY
					WIN_RATIO DESC,
					NAME ASC
				LIMIT 5
			`
		default:
			return resultHeaders, resultRows, errors.New("invalid subject provided")
		}
	case "game-win":
		switch subject {
		case "player":
			resultHeaders = append(resultHeaders, "Games Won")
			resultHeaders = append(resultHeaders, "Player")
			sqlString = `
				SELECT
					COUNT(*) AS COUNT,
					U.NAME AS NAME
				FROM (
						SELECT
							LRC.LOBBY_ID,
							LRC.PLAYER_USER_ID,
							COUNT(LW.ID) AS ROUND_WIN_COUNT,
							RANK() OVER (PARTITION BY LRC.LOBBY_ID ORDER BY ROUND_WIN_COUNT DESC) AS RANKING
						FROM LOG_RESPONSE_CARD AS LRC
							INNER JOIN LOG_WIN AS LW ON LW.RESPONSE_ID = LRC.RESPONSE_ID
						GROUP BY LRC.LOBBY_ID, LRC.PLAYER_USER_ID
					) AS RW
					INNER JOIN USER AS U ON U.ID = RW.PLAYER_USER_ID
				WHERE RW.RANKING = 1
				GROUP BY RW.PLAYER_USER_ID
				ORDER BY
					COUNT DESC,
					NAME ASC
				LIMIT 5
			`
		default:
			return resultHeaders, resultRows, errors.New("invalid subject provided")
		}
	case "game-play":
		switch subject {
		case "player":
			resultHeaders = append(resultHeaders, "Games Played")
			resultHeaders = append(resultHeaders, "Player")
			sqlString = `
				SELECT
					COUNT(DISTINCT LRC.LOBBY_ID) AS COUNT,
					U.NAME AS NAME
				FROM LOG_RESPONSE_CARD AS LRC
					INNER JOIN USER AS U ON U.ID = LRC.PLAYER_USER_ID
				GROUP BY U.ID
				ORDER BY
					COUNT DESC,
					NAME ASC
				LIMIT 5
			`
		case "card":
			resultHeaders = append(resultHeaders, "Games Played")
			resultHeaders = append(resultHeaders, "Card")
			sqlString = `
				SELECT
					COUNT(DISTINCT LRC.LOBBY_ID) AS COUNT,
					COALESCE(C.TEXT, LRC.SPECIAL_CATEGORY, 'Unknown') AS NAME
				FROM LOG_RESPONSE_CARD AS LRC
					LEFT JOIN CARD AS C ON C.ID = LRC.PLAYER_CARD_ID
				GROUP BY C.ID
				ORDER BY
					COUNT DESC,
					NAME ASC
				LIMIT 5
			`
		case "special-category":
			resultHeaders = append(resultHeaders, "Games Played")
			resultHeaders = append(resultHeaders, "Special Category")
			sqlString = `
				SELECT
					COUNT(DISTINCT LRC.LOBBY_ID) AS COUNT,
					COALESCE(LRC.SPECIAL_CATEGORY, 'NONE') AS NAME
				FROM LOG_RESPONSE_CARD AS LRC
				GROUP BY LRC.SPECIAL_CATEGORY
				ORDER BY
					COUNT DESC,
					NAME ASC
				LIMIT 5
			`
		default:
			return resultHeaders, resultRows, errors.New("invalid subject provided")
		}
	case "round-win-ratio":
		switch subject {
		case "player":
			resultHeaders = append(resultHeaders, "Rounds Played")
			resultHeaders = append(resultHeaders, "Rounds Won")
			resultHeaders = append(resultHeaders, "Win Ratio")
			resultHeaders = append(resultHeaders, "Player")
			sqlString = `
				SELECT
					PLAY_COUNT,
					WIN_COUNT,
					COALESCE((WIN_COUNT * 1.0) / (PLAY_COUNT * 1.0), 0.0) AS WIN_RATIO,
					NAME
				FROM (
						SELECT
							COUNT(DISTINCT LRC.ROUND_ID) AS PLAY_COUNT,
							COUNT(DISTINCT LW.ID)        AS WIN_COUNT,
							U.NAME                       AS NAME
						FROM LOG_RESPONSE_CARD AS LRC
							LEFT JOIN LOG_WIN AS LW ON LW.RESPONSE_ID = LRC.RESPONSE_ID
							INNER JOIN USER AS U ON U.ID = LRC.PLAYER_USER_ID
						GROUP BY U.ID
					) AS T
				ORDER BY
					WIN_RATIO DESC,
					NAME ASC
				LIMIT 5
			`
		case "card":
			resultHeaders = append(resultHeaders, "Rounds Played")
			resultHeaders = append(resultHeaders, "Rounds Won")
			resultHeaders = append(resultHeaders, "Win Ratio")
			resultHeaders = append(resultHeaders, "Card")
			params = append(params, userId)
			sqlString = `
				SELECT
					PLAY_COUNT,
					WIN_COUNT,
					COALESCE((WIN_COUNT * 1.0) / (PLAY_COUNT * 1.0), 0.0) AS WIN_RATIO,
					NAME
				FROM (
						SELECT
							COUNT(DISTINCT LRC.ROUND_ID) AS PLAY_COUNT,
							COUNT(DISTINCT LW.ID)        AS WIN_COUNT,
							COALESCE(C.TEXT, LRC.SPECIAL_CATEGORY, 'Unknown') AS NAME
						FROM LOG_RESPONSE_CARD AS LRC
							LEFT JOIN LOG_WIN AS LW ON LW.RESPONSE_ID = LRC.RESPONSE_ID
							LEFT JOIN CARD AS C ON C.ID = LRC.PLAYER_CARD_ID
						WHERE FN_USER_HAS_DECK_ACCESS(?, C.DECK_ID)
						GROUP BY C.ID
					) AS T
				ORDER BY
					WIN_RATIO DESC,
					NAME ASC
				LIMIT 5
			`
		case "special-category":
			resultHeaders = append(resultHeaders, "Rounds Played")
			resultHeaders = append(resultHeaders, "Rounds Won")
			resultHeaders = append(resultHeaders, "Win Ratio")
			resultHeaders = append(resultHeaders, "Special Category")
			sqlString = `
				SELECT
					PLAY_COUNT,
					WIN_COUNT,
					COALESCE((WIN_COUNT * 1.0) / (PLAY_COUNT * 1.0), 0.0) AS WIN_RATIO,
					NAME
				FROM (
						SELECT
							COUNT(DISTINCT LRC.ROUND_ID) AS PLAY_COUNT,
							COUNT(DISTINCT LW.ID)        AS WIN_COUNT,
							COALESCE(LRC.SPECIAL_CATEGORY, 'NONE') AS NAME
						FROM LOG_RESPONSE_CARD AS LRC
							LEFT JOIN LOG_WIN AS LW ON LW.RESPONSE_ID = LRC.RESPONSE_ID
						GROUP BY LRC.SPECIAL_CATEGORY
					) AS T
				ORDER BY
					WIN_RATIO DESC,
					NAME ASC
				LIMIT 5
			`
		default:
			return resultHeaders, resultRows, errors.New("invalid subject provided")
		}
	case "round-win":
		switch subject {
		case "player":
			resultHeaders = append(resultHeaders, "Rounds Won")
			resultHeaders = append(resultHeaders, "Player")
			sqlString = `
				SELECT
					COUNT(DISTINCT LRC.ROUND_ID) AS COUNT,
					U.NAME AS NAME
				FROM LOG_RESPONSE_CARD AS LRC
					INNER JOIN LOG_WIN AS LW ON LW.RESPONSE_ID = LRC.RESPONSE_ID
					INNER JOIN USER AS U ON U.ID = LRC.PLAYER_USER_ID
				GROUP BY U.ID
				ORDER BY
					COUNT DESC,
					NAME ASC
				LIMIT 5
			`
		case "card":
			resultHeaders = append(resultHeaders, "Rounds Won")
			resultHeaders = append(resultHeaders, "Card")
			params = append(params, userId)
			sqlString = `
				SELECT
					COUNT(DISTINCT LRC.ROUND_ID) AS COUNT,
					COALESCE(C.TEXT, LRC.SPECIAL_CATEGORY, 'Unknown') AS NAME
				FROM LOG_RESPONSE_CARD AS LRC
					INNER JOIN LOG_WIN AS LW ON LW.RESPONSE_ID = LRC.RESPONSE_ID
					LEFT JOIN CARD AS C ON C.ID = LRC.PLAYER_CARD_ID
				WHERE FN_USER_HAS_DECK_ACCESS(?, C.DECK_ID)
				GROUP BY C.ID
				ORDER BY
					COUNT DESC,
					NAME ASC
				LIMIT 5
			`
		case "special-category":
			resultHeaders = append(resultHeaders, "Rounds Won")
			resultHeaders = append(resultHeaders, "Special Category")
			sqlString = `
				SELECT
					COUNT(DISTINCT LRC.ROUND_ID) AS COUNT,
					COALESCE(LRC.SPECIAL_CATEGORY, 'NONE') AS NAME
				FROM LOG_RESPONSE_CARD AS LRC
					INNER JOIN LOG_WIN AS LW ON LW.RESPONSE_ID = LRC.RESPONSE_ID
				GROUP BY LRC.SPECIAL_CATEGORY
				ORDER BY
					COUNT DESC,
					NAME ASC
				LIMIT 5
			`
		default:
			return resultHeaders, resultRows, errors.New("invalid subject provided")
		}
	case "round-play":
		switch subject {
		case "player":
			resultHeaders = append(resultHeaders, "Rounds Played")
			resultHeaders = append(resultHeaders, "Player")
			sqlString = `
				SELECT
					COUNT(DISTINCT LRC.ROUND_ID) AS COUNT,
					U.NAME AS NAME
				FROM LOG_RESPONSE_CARD AS LRC
					INNER JOIN USER AS U ON U.ID = LRC.PLAYER_USER_ID
				GROUP BY U.ID
				ORDER BY
					COUNT DESC,
					NAME ASC
				LIMIT 5
			`
		case "card":
			resultHeaders = append(resultHeaders, "Rounds Played")
			resultHeaders = append(resultHeaders, "Card")
			params = append(params, userId)
			sqlString = `
				SELECT
					COUNT(DISTINCT LRC.ROUND_ID) AS COUNT,
					COALESCE(C.TEXT, LRC.SPECIAL_CATEGORY, 'Unknown') AS NAME
				FROM LOG_RESPONSE_CARD AS LRC
					LEFT JOIN CARD AS C ON C.ID = LRC.PLAYER_CARD_ID
				WHERE FN_USER_HAS_DECK_ACCESS(?, C.DECK_ID)
				GROUP BY C.ID
				ORDER BY
					COUNT DESC,
					NAME ASC
				LIMIT 5
			`
		case "special-category":
			resultHeaders = append(resultHeaders, "Rounds Played")
			resultHeaders = append(resultHeaders, "Special Category")
			sqlString = `
				SELECT
					COUNT(DISTINCT LRC.ROUND_ID) AS COUNT,
					COALESCE(LRC.SPECIAL_CATEGORY, 'NONE') AS NAME
				FROM LOG_RESPONSE_CARD AS LRC
				GROUP BY LRC.SPECIAL_CATEGORY
				ORDER BY
					COUNT DESC,
					NAME ASC
				LIMIT 5
			`
		default:
			return resultHeaders, resultRows, errors.New("invalid subject provided")
		}
	case "card-play":
		switch subject {
		case "player":
			resultHeaders = append(resultHeaders, "Cards Played")
			resultHeaders = append(resultHeaders, "Player")
			sqlString = `
				SELECT
					COUNT(DISTINCT LRC.ID) AS COUNT,
					U.NAME AS NAME
				FROM LOG_RESPONSE_CARD AS LRC
					INNER JOIN USER AS U ON U.ID = LRC.PLAYER_USER_ID
				GROUP BY U.ID
				ORDER BY
					COUNT DESC,
					NAME ASC
				LIMIT 5
			`
		case "card":
			resultHeaders = append(resultHeaders, "Cards Played")
			resultHeaders = append(resultHeaders, "Card")
			params = append(params, userId)
			sqlString = `
				SELECT
					COUNT(DISTINCT LRC.ID) AS COUNT,
					COALESCE(C.TEXT, LRC.SPECIAL_CATEGORY, 'Unknown') AS NAME
				FROM LOG_RESPONSE_CARD AS LRC
					LEFT JOIN CARD AS C ON C.ID = LRC.PLAYER_CARD_ID
				WHERE FN_USER_HAS_DECK_ACCESS(?, C.DECK_ID)
				GROUP BY C.ID
				ORDER BY
					COUNT DESC,
					NAME ASC
				LIMIT 5
			`
		case "special-category":
			resultHeaders = append(resultHeaders, "Cards Played")
			resultHeaders = append(resultHeaders, "Special Category")
			sqlString = `
				SELECT
					COUNT(DISTINCT LRC.ID) AS COUNT,
					COALESCE(LRC.SPECIAL_CATEGORY, 'NONE') AS NAME
				FROM LOG_RESPONSE_CARD AS LRC
				GROUP BY LRC.SPECIAL_CATEGORY
				ORDER BY
					COUNT DESC,
					NAME ASC
				LIMIT 5
			`
		default:
			return resultHeaders, resultRows, errors.New("invalid subject provided")
		}
	case "card-draw":
		switch subject {
		case "player":
			resultHeaders = append(resultHeaders, "Cards Drawn")
			resultHeaders = append(resultHeaders, "Player")
			sqlString = `
				SELECT
					COUNT(DISTINCT LD.ID) AS COUNT,
					U.NAME AS NAME
				FROM LOG_DRAW AS LD
					INNER JOIN USER AS U ON U.ID = LD.USER_ID
				GROUP BY U.ID
				ORDER BY
					COUNT DESC,
					NAME ASC
				LIMIT 5
			`
		case "card":
			resultHeaders = append(resultHeaders, "Cards Drawn")
			resultHeaders = append(resultHeaders, "Card")
			params = append(params, userId)
			sqlString = `
				SELECT
					COUNT(DISTINCT LD.ID) AS COUNT,
					COALESCE(C.TEXT, 'Unknown') AS NAME
				FROM LOG_DRAW AS LD
					LEFT JOIN CARD AS C ON C.ID = LD.CARD_ID
				WHERE FN_USER_HAS_DECK_ACCESS(?, C.DECK_ID)
				GROUP BY C.ID
				ORDER BY
					COUNT DESC,
					NAME ASC
				LIMIT 5
			`
		default:
			return resultHeaders, resultRows, errors.New("invalid subject provided")
		}
	case "card-discard":
		switch subject {
		case "player":
			resultHeaders = append(resultHeaders, "Cards Discarded")
			resultHeaders = append(resultHeaders, "Player")
			sqlString = `
				SELECT
					COUNT(DISTINCT LD.ID) AS COUNT,
					U.NAME AS NAME
				FROM LOG_DISCARD AS LD
					INNER JOIN USER AS U ON U.ID = LD.USER_ID
				GROUP BY U.ID
				ORDER BY
					COUNT DESC,
					NAME ASC
				LIMIT 5
			`
		case "card":
			resultHeaders = append(resultHeaders, "Cards Discarded")
			resultHeaders = append(resultHeaders, "Card")
			params = append(params, userId)
			sqlString = `
				SELECT
					COUNT(DISTINCT LD.ID) AS COUNT,
					COALESCE(C.TEXT, 'Unknown') AS NAME
				FROM LOG_DISCARD AS LD
					LEFT JOIN CARD AS C ON C.ID = LD.CARD_ID
				WHERE FN_USER_HAS_DECK_ACCESS(?, C.DECK_ID)
				GROUP BY C.ID
				ORDER BY
					COUNT DESC,
					NAME ASC
				LIMIT 5
			`
		default:
			return resultHeaders, resultRows, errors.New("invalid subject provided")
		}
	case "card-skip":
		switch subject {
		case "player":
			resultHeaders = append(resultHeaders, "Cards Skipped")
			resultHeaders = append(resultHeaders, "Player")
			sqlString = `
				SELECT
					COUNT(DISTINCT LS.ID) AS COUNT,
					U.NAME AS NAME
				FROM LOG_SKIP AS LS
					INNER JOIN USER AS U ON U.ID = LS.USER_ID
				GROUP BY U.ID
				ORDER BY
					COUNT DESC,
					NAME ASC
				LIMIT 5
			`
		case "card":
			resultHeaders = append(resultHeaders, "Cards Skipped")
			resultHeaders = append(resultHeaders, "Card")
			params = append(params, userId)
			sqlString = `
				SELECT
					COUNT(DISTINCT LS.ID) AS COUNT,
					COALESCE(C.TEXT, 'Unknown') AS NAME
				FROM LOG_SKIP AS LS
					LEFT JOIN CARD AS C ON C.ID = LS.CARD_ID
				WHERE FN_USER_HAS_DECK_ACCESS(?, C.DECK_ID)
				GROUP BY C.ID
				ORDER BY
					COUNT DESC,
					NAME ASC
				LIMIT 5
			`
		default:
			return resultHeaders, resultRows, errors.New("invalid subject provided")
		}
	case "picked-judge":
		switch subject {
		case "player":
			resultHeaders = append(resultHeaders, "Judge Picking")
			resultHeaders = append(resultHeaders, "Player")
			resultHeaders = append(resultHeaders, "Count")
			params = append(params, userId)
			sqlString = `
				SELECT
					UJ.NAME AS JUDGE_NAME,
					UP.NAME AS NAME,
					COUNT(LW.ID) AS COUNT
				FROM LOG_RESPONSE_CARD AS LRC
					INNER JOIN LOG_WIN AS LW ON LW.RESPONSE_ID = LRC.RESPONSE_ID
					INNER JOIN USER AS UJ ON UJ.ID = LRC.JUDGE_USER_ID
					INNER JOIN USER AS UP ON UP.ID = LRC.PLAYER_USER_ID
				WHERE UJ.ID = ?
				GROUP BY LRC.JUDGE_USER_ID, LRC.PLAYER_USER_ID
				ORDER BY
					COUNT DESC,
					NAME ASC
				LIMIT 5
			`
		case "card":
			resultHeaders = append(resultHeaders, "Judge Picking")
			resultHeaders = append(resultHeaders, "Card")
			resultHeaders = append(resultHeaders, "Count")
			params = append(params, userId)
			sqlString = `
				SELECT
					UJ.NAME AS JUDGE_NAME,
					CP.TEXT AS NAME,
					COUNT(LW.ID) AS COUNT
				FROM LOG_RESPONSE_CARD AS LRC
					INNER JOIN LOG_WIN AS LW ON LW.RESPONSE_ID = LRC.RESPONSE_ID
					INNER JOIN USER AS UJ ON UJ.ID = LRC.JUDGE_USER_ID
					INNER JOIN CARD AS CP ON CP.ID = LRC.PLAYER_CARD_ID
				WHERE UJ.ID = ?
				GROUP BY LRC.JUDGE_USER_ID, LRC.PLAYER_CARD_ID
				ORDER BY
					COUNT DESC,
					NAME ASC
				LIMIT 5
			`
		case "special-category":
			resultHeaders = append(resultHeaders, "Judge Picking")
			resultHeaders = append(resultHeaders, "Special Category")
			resultHeaders = append(resultHeaders, "Count")
			params = append(params, userId)
			sqlString = `
				SELECT
					UJ.NAME AS JUDGE_NAME,
					COALESCE(LRC.SPECIAL_CATEGORY, 'NONE') AS NAME,
					COUNT(LW.ID) AS COUNT
				FROM LOG_RESPONSE_CARD AS LRC
					INNER JOIN LOG_WIN AS LW ON LW.RESPONSE_ID = LRC.RESPONSE_ID
					INNER JOIN USER AS UJ ON UJ.ID = LRC.JUDGE_USER_ID
				WHERE UJ.ID = ?
				GROUP BY LRC.JUDGE_USER_ID, LRC.SPECIAL_CATEGORY
				ORDER BY
					COUNT DESC,
					NAME ASC
				LIMIT 5
			`
		default:
			return resultHeaders, resultRows, errors.New("invalid subject provided")
		}
	case "picked-player":
		switch subject {
		case "player":
			resultHeaders = append(resultHeaders, "Winner")
			resultHeaders = append(resultHeaders, "Judge Who Picked")
			resultHeaders = append(resultHeaders, "Count")
			params = append(params, userId)
			sqlString = `
				SELECT
					UP.NAME AS PLAYER_NAME,
					UJ.NAME AS NAME,
					COUNT(LW.ID) AS COUNT
				FROM LOG_RESPONSE_CARD AS LRC
					INNER JOIN LOG_WIN AS LW ON LW.RESPONSE_ID = LRC.RESPONSE_ID
					INNER JOIN USER AS UJ ON UJ.ID = LRC.JUDGE_USER_ID
					INNER JOIN USER AS UP ON UP.ID = LRC.PLAYER_USER_ID
				WHERE UP.ID = ?
				GROUP BY LRC.JUDGE_USER_ID, LRC.PLAYER_USER_ID
				ORDER BY
					COUNT DESC,
					NAME ASC
				LIMIT 5
			`
		case "card":
			resultHeaders = append(resultHeaders, "Winner")
			resultHeaders = append(resultHeaders, "Card Played")
			resultHeaders = append(resultHeaders, "Count")
			params = append(params, userId)
			sqlString = `
				SELECT
					UP.NAME AS PLAYER_NAME,
					CJ.TEXT AS NAME,
					COUNT(LW.ID) AS COUNT
				FROM LOG_RESPONSE_CARD AS LRC
					INNER JOIN LOG_WIN AS LW ON LW.RESPONSE_ID = LRC.RESPONSE_ID
					INNER JOIN CARD AS CJ ON CJ.ID = LRC.PLAYER_CARD_ID
					INNER JOIN USER AS UP ON UP.ID = LRC.PLAYER_USER_ID
				WHERE UP.ID = ?
				GROUP BY LRC.PLAYER_CARD_ID, LRC.PLAYER_USER_ID
				ORDER BY
					COUNT DESC,
					NAME ASC
				LIMIT 5
			`
		case "special-category":
			resultHeaders = append(resultHeaders, "Winner")
			resultHeaders = append(resultHeaders, "Special Category Played")
			resultHeaders = append(resultHeaders, "Count")
			params = append(params, userId)
			sqlString = `
				SELECT
					UP.NAME AS PLAYER_NAME,
					COALESCE(LRC.SPECIAL_CATEGORY, 'NONE') AS NAME,
					COUNT(LW.ID) AS COUNT
				FROM LOG_RESPONSE_CARD AS LRC
					INNER JOIN LOG_WIN AS LW ON LW.RESPONSE_ID = LRC.RESPONSE_ID
					INNER JOIN USER AS UP ON UP.ID = LRC.PLAYER_USER_ID
				WHERE UP.ID = ?
				GROUP BY LRC.SPECIAL_CATEGORY, LRC.PLAYER_USER_ID
				ORDER BY
					COUNT DESC,
					NAME ASC
				LIMIT 5
			`
		default:
			return resultHeaders, resultRows, errors.New("invalid subject provided")
		}
	default:
		return resultHeaders, resultRows, errors.New("invalid topic provided")
	}

	rows, err := query(sqlString, params...)
	if err != nil {
		return resultHeaders, resultRows, err
	}

	rowValuePointers := make([]any, len(resultHeaders))
	for i := range rowValuePointers {
		rowValuePointers[i] = new(string)
	}

	for rows.Next() {
		if err := rows.Scan(rowValuePointers...); err != nil {
			log.Println(err)
			return resultHeaders, resultRows, errors.New("failed to scan row in query results")
		}

		row := make([]string, len(resultHeaders))
		for i, v := range rowValuePointers {
			row[i] = reflect.ValueOf(v).Elem().String()
		}
		resultRows = append(resultRows, row)
	}

	return resultHeaders, resultRows, nil
}
