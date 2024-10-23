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

func GetMostWinsByPlayer() ([]MostWins, error) {
	sqlString := `
		SELECT
			COUNT(LW.ID) AS WIN_COUNT,
			U.NAME AS NAME
		FROM LOG_WIN AS LW
			INNER JOIN USER AS U ON U.ID = LW.PLAYER_USER_ID
		GROUP BY U.ID
		ORDER BY
			COUNT(LW.ID) DESC,
			U.NAME ASC
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
			COUNT(LW.ID) AS WIN_COUNT,
			COALESCE(C.TEXT, LW.SPECIAL_CATEGORY, 'Unknown') AS NAME
		FROM LOG_WIN AS LW
			LEFT JOIN CARD AS C ON C.ID = LW.CARD_ID
		WHERE FN_USER_HAS_DECK_ACCESS(?, C.DECK_ID)
		GROUP BY C.ID
		ORDER BY
			COUNT(LW.ID) DESC,
			C.TEXT ASC
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
