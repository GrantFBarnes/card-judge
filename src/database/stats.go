package database

import (
	"errors"
	"log"

	"github.com/google/uuid"
)

type MostWins struct {
	WinCount int
	Name     string
}

type BestWinRatio struct {
	PlayCount int
	WinCount  int
	WinRatio  float64
	Name      string
}

func GetMostWinsByPlayer() ([]MostWins, error) {
	sqlString := `
		SELECT
			COUNT(DISTINCT LW.ID) AS WIN_COUNT,
			U.NAME AS NAME
		FROM LOG_WIN AS LW
			INNER JOIN USER AS U ON U.ID = LW.PLAYER_USER_ID
		GROUP BY U.ID
		ORDER BY
			WIN_COUNT DESC,
			NAME ASC
		LIMIT 5
	`
	rows, err := query(sqlString)
	if err != nil {
		return nil, err
	}

	result := make([]MostWins, 0)
	for rows.Next() {
		var mw MostWins
		if err := rows.Scan(&mw.WinCount, &mw.Name); err != nil {
			log.Println(err)
			return result, errors.New("failed to scan row in query results")
		}
		result = append(result, mw)
	}

	return result, nil
}

func GetMostWinsByCard(userId uuid.UUID) ([]MostWins, error) {
	sqlString := `
		SELECT
			COUNT(DISTINCT LW.ID) AS WIN_COUNT,
			COALESCE(C.TEXT, LW.SPECIAL_CATEGORY, 'Unknown') AS NAME
		FROM LOG_WIN AS LW
			LEFT JOIN CARD AS C ON C.ID = LW.CARD_ID
		WHERE FN_USER_HAS_DECK_ACCESS(?, C.DECK_ID)
		GROUP BY C.ID
		ORDER BY
			WIN_COUNT DESC,
			NAME ASC
		LIMIT 5
	`
	rows, err := query(sqlString, userId)
	if err != nil {
		return nil, err
	}

	result := make([]MostWins, 0)
	for rows.Next() {
		var mw MostWins
		if err := rows.Scan(&mw.WinCount, &mw.Name); err != nil {
			log.Println(err)
			return result, errors.New("failed to scan row in query results")
		}
		result = append(result, mw)
	}

	return result, nil
}

func GetMostWinsBySpecialCategory() ([]MostWins, error) {
	sqlString := `
		SELECT
			COUNT(DISTINCT LW.ID) AS WIN_COUNT,
			COALESCE(LW.SPECIAL_CATEGORY, 'NONE') AS NAME
		FROM LOG_WIN AS LW
		GROUP BY LW.SPECIAL_CATEGORY
		ORDER BY
			WIN_COUNT DESC,
			NAME ASC
		LIMIT 5
	`
	rows, err := query(sqlString)
	if err != nil {
		return nil, err
	}

	result := make([]MostWins, 0)
	for rows.Next() {
		var mw MostWins
		if err := rows.Scan(&mw.WinCount, &mw.Name); err != nil {
			log.Println(err)
			return result, errors.New("failed to scan row in query results")
		}
		result = append(result, mw)
	}

	return result, nil
}

func GetBestWinRatioByPlayer() ([]BestWinRatio, error) {
	sqlString := `
		SELECT
			COUNT(DISTINCT LP.ID) AS PLAY_COUNT,
			COUNT(DISTINCT LW.ID) AS WIN_COUNT,
			(COUNT(DISTINCT LW.ID)*1.0) / (COUNT(DISTINCT LP.ID)*1.0) AS WIN_RATIO,
			U.NAME AS NAME
		FROM LOG_WIN AS LW
			INNER JOIN USER AS U ON U.ID = LW.PLAYER_USER_ID
			INNER JOIN LOG_PLAY AS LP ON LP.PLAYER_USER_ID = U.ID
		GROUP BY U.ID
		ORDER BY
			WIN_RATIO DESC,
			NAME ASC
		LIMIT 5
	`
	rows, err := query(sqlString)
	if err != nil {
		return nil, err
	}

	result := make([]BestWinRatio, 0)
	for rows.Next() {
		var bwr BestWinRatio
		if err := rows.Scan(
			&bwr.PlayCount,
			&bwr.WinCount,
			&bwr.WinRatio,
			&bwr.Name); err != nil {
			log.Println(err)
			return result, errors.New("failed to scan row in query results")
		}
		result = append(result, bwr)
	}

	return result, nil
}

func GetBestWinRatioByCard(userId uuid.UUID) ([]BestWinRatio, error) {
	sqlString := `
		SELECT
			COUNT(DISTINCT LP.ID) AS PLAY_COUNT,
			COUNT(DISTINCT LW.ID) AS WIN_COUNT,
			(COUNT(DISTINCT LW.ID)*1.0) / (COUNT(DISTINCT LP.ID)*1.0) AS WIN_RATIO,
			COALESCE(C.TEXT, LW.SPECIAL_CATEGORY, 'Unknown') AS NAME
		FROM LOG_WIN AS LW
			LEFT JOIN CARD AS C ON C.ID = LW.CARD_ID
			LEFT JOIN LOG_PLAY AS LP ON LP.CARD_ID = C.ID
		WHERE FN_USER_HAS_DECK_ACCESS(?, C.DECK_ID)
		GROUP BY C.ID
		ORDER BY
			WIN_RATIO DESC,
			NAME ASC
		LIMIT 5
	`
	rows, err := query(sqlString, userId)
	if err != nil {
		return nil, err
	}

	result := make([]BestWinRatio, 0)
	for rows.Next() {
		var bwr BestWinRatio
		if err := rows.Scan(
			&bwr.PlayCount,
			&bwr.WinCount,
			&bwr.WinRatio,
			&bwr.Name); err != nil {
			log.Println(err)
			return result, errors.New("failed to scan row in query results")
		}
		result = append(result, bwr)
	}

	return result, nil
}

func GetBestWinRatioBySpecialCategory() ([]BestWinRatio, error) {
	sqlString := `
		SELECT
			COUNT(DISTINCT LP.ID) AS PLAY_COUNT,
			COUNT(DISTINCT LW.ID) AS WIN_COUNT,
			(COUNT(DISTINCT LW.ID)*1.0) / (COUNT(DISTINCT LP.ID)*1.0) AS WIN_RATIO,
			COALESCE(LW.SPECIAL_CATEGORY, 'NONE') AS NAME
		FROM LOG_WIN AS LW
			LEFT JOIN LOG_PLAY AS LP ON LP.CARD_ID = LW.CARD_ID OR LP.CARD_ID IS NULL
		GROUP BY LW.SPECIAL_CATEGORY
		ORDER BY
			WIN_RATIO DESC,
			NAME ASC
		LIMIT 5
	`
	rows, err := query(sqlString)
	if err != nil {
		return nil, err
	}

	result := make([]BestWinRatio, 0)
	for rows.Next() {
		var bwr BestWinRatio
		if err := rows.Scan(
			&bwr.PlayCount,
			&bwr.WinCount,
			&bwr.WinRatio,
			&bwr.Name); err != nil {
			log.Println(err)
			return result, errors.New("failed to scan row in query results")
		}
		result = append(result, bwr)
	}

	return result, nil
}
