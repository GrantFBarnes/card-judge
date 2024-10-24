package database

import (
	"errors"
	"log"

	"github.com/google/uuid"
)

type StatCount struct {
	Count int
	Name  string
}

type StatWinRatio struct {
	PlayCount int
	WinCount  int
	WinRatio  float64
	Name      string
}

func GetMostPlaysByPlayer() ([]StatCount, error) {
	sqlString := `
		SELECT
			COUNT(DISTINCT LP.ID) AS COUNT,
			U.NAME AS NAME
		FROM LOG_PLAY AS LP
			INNER JOIN USER AS U ON U.ID = LP.PLAYER_USER_ID
		GROUP BY U.ID
		ORDER BY
			COUNT DESC,
			NAME ASC
		LIMIT 5
	`
	rows, err := query(sqlString)
	if err != nil {
		return nil, err
	}

	result := make([]StatCount, 0)
	for rows.Next() {
		var sc StatCount
		if err := rows.Scan(&sc.Count, &sc.Name); err != nil {
			log.Println(err)
			return result, errors.New("failed to scan row in query results")
		}
		result = append(result, sc)
	}

	return result, nil
}

func GetMostPlaysByCard(userId uuid.UUID) ([]StatCount, error) {
	sqlString := `
		SELECT
			COUNT(DISTINCT LP.ID) AS COUNT,
			COALESCE(C.TEXT, LP.SPECIAL_CATEGORY, 'UNKNOWN') AS NAME
		FROM LOG_PLAY AS LP
			LEFT JOIN CARD AS C ON C.ID = LP.CARD_ID
		WHERE FN_USER_HAS_DECK_ACCESS(?, C.DECK_ID)
		GROUP BY C.ID
		ORDER BY
			COUNT DESC,
			NAME ASC
		LIMIT 5
	`
	rows, err := query(sqlString, userId)
	if err != nil {
		return nil, err
	}

	result := make([]StatCount, 0)
	for rows.Next() {
		var sc StatCount
		if err := rows.Scan(&sc.Count, &sc.Name); err != nil {
			log.Println(err)
			return result, errors.New("failed to scan row in query results")
		}
		result = append(result, sc)
	}

	return result, nil
}

func GetMostPlaysBySpecialCategory() ([]StatCount, error) {
	sqlString := `
		SELECT
			COUNT(DISTINCT LP.ID) AS COUNT,
			COALESCE(LP.SPECIAL_CATEGORY, 'NONE') AS NAME
		FROM LOG_PLAY AS LP
		GROUP BY LP.SPECIAL_CATEGORY
		ORDER BY
			COUNT DESC,
			NAME ASC
		LIMIT 5
	`
	rows, err := query(sqlString)
	if err != nil {
		return nil, err
	}

	result := make([]StatCount, 0)
	for rows.Next() {
		var sc StatCount
		if err := rows.Scan(&sc.Count, &sc.Name); err != nil {
			log.Println(err)
			return result, errors.New("failed to scan row in query results")
		}
		result = append(result, sc)
	}

	return result, nil
}

func GetMostWinsByPlayer() ([]StatCount, error) {
	sqlString := `
		SELECT
			COUNT(DISTINCT LW.ID) AS COUNT,
			U.NAME AS NAME
		FROM LOG_WIN AS LW
			INNER JOIN USER AS U ON U.ID = LW.PLAYER_USER_ID
		GROUP BY U.ID
		ORDER BY
			COUNT DESC,
			NAME ASC
		LIMIT 5
	`
	rows, err := query(sqlString)
	if err != nil {
		return nil, err
	}

	result := make([]StatCount, 0)
	for rows.Next() {
		var sc StatCount
		if err := rows.Scan(&sc.Count, &sc.Name); err != nil {
			log.Println(err)
			return result, errors.New("failed to scan row in query results")
		}
		result = append(result, sc)
	}

	return result, nil
}

func GetMostWinsByCard(userId uuid.UUID) ([]StatCount, error) {
	sqlString := `
		SELECT
			COUNT(DISTINCT LW.ID) AS COUNT,
			COALESCE(C.TEXT, LW.SPECIAL_CATEGORY, 'Unknown') AS NAME
		FROM LOG_WIN AS LW
			LEFT JOIN CARD AS C ON C.ID = LW.CARD_ID
		WHERE FN_USER_HAS_DECK_ACCESS(?, C.DECK_ID)
		GROUP BY C.ID
		ORDER BY
			COUNT DESC,
			NAME ASC
		LIMIT 5
	`
	rows, err := query(sqlString, userId)
	if err != nil {
		return nil, err
	}

	result := make([]StatCount, 0)
	for rows.Next() {
		var sc StatCount
		if err := rows.Scan(&sc.Count, &sc.Name); err != nil {
			log.Println(err)
			return result, errors.New("failed to scan row in query results")
		}
		result = append(result, sc)
	}

	return result, nil
}

func GetMostWinsBySpecialCategory() ([]StatCount, error) {
	sqlString := `
		SELECT
			COUNT(DISTINCT LW.ID) AS COUNT,
			COALESCE(LW.SPECIAL_CATEGORY, 'NONE') AS NAME
		FROM LOG_WIN AS LW
		GROUP BY LW.SPECIAL_CATEGORY
		ORDER BY
			COUNT DESC,
			NAME ASC
		LIMIT 5
	`
	rows, err := query(sqlString)
	if err != nil {
		return nil, err
	}

	result := make([]StatCount, 0)
	for rows.Next() {
		var sc StatCount
		if err := rows.Scan(&sc.Count, &sc.Name); err != nil {
			log.Println(err)
			return result, errors.New("failed to scan row in query results")
		}
		result = append(result, sc)
	}

	return result, nil
}

func GetBestWinRatioByPlayer() ([]StatWinRatio, error) {
	sqlString := `
		SELECT
			COUNT(DISTINCT LP.ID) AS PLAY_COUNT,
			COUNT(DISTINCT LW.ID) AS WIN_COUNT,
			COALESCE((COUNT(DISTINCT LW.ID)*1.0) / (COUNT(DISTINCT LP.ID)*1.0),0.0) AS WIN_RATIO,
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

	result := make([]StatWinRatio, 0)
	for rows.Next() {
		var swr StatWinRatio
		if err := rows.Scan(
			&swr.PlayCount,
			&swr.WinCount,
			&swr.WinRatio,
			&swr.Name); err != nil {
			log.Println(err)
			return result, errors.New("failed to scan row in query results")
		}
		result = append(result, swr)
	}

	return result, nil
}

func GetBestWinRatioByCard(userId uuid.UUID) ([]StatWinRatio, error) {
	sqlString := `
		SELECT
			COUNT(DISTINCT LP.ID) AS PLAY_COUNT,
			COUNT(DISTINCT LW.ID) AS WIN_COUNT,
			COALESCE((COUNT(DISTINCT LW.ID)*1.0) / (COUNT(DISTINCT LP.ID)*1.0),0.0) AS WIN_RATIO,
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

	result := make([]StatWinRatio, 0)
	for rows.Next() {
		var swr StatWinRatio
		if err := rows.Scan(
			&swr.PlayCount,
			&swr.WinCount,
			&swr.WinRatio,
			&swr.Name); err != nil {
			log.Println(err)
			return result, errors.New("failed to scan row in query results")
		}
		result = append(result, swr)
	}

	return result, nil
}

func GetBestWinRatioBySpecialCategory() ([]StatWinRatio, error) {
	sqlString := `
		SELECT
			COUNT(DISTINCT LP.ID) AS PLAY_COUNT,
			COUNT(DISTINCT LW.ID) AS WIN_COUNT,
			COALESCE((COUNT(DISTINCT LW.ID)*1.0) / (COUNT(DISTINCT LP.ID)*1.0),0.0) AS WIN_RATIO,
			COALESCE(LP.SPECIAL_CATEGORY, 'NONE') AS NAME
		FROM LOG_PLAY AS LP
			INNER JOIN LOG_WIN AS LW ON (LW.SPECIAL_CATEGORY = LP.SPECIAL_CATEGORY) OR
										(LW.SPECIAL_CATEGORY IS NULL AND LP.SPECIAL_CATEGORY IS NULL)
		GROUP BY LP.SPECIAL_CATEGORY
		ORDER BY
			WIN_RATIO DESC,
			NAME ASC
		LIMIT 5
	`
	rows, err := query(sqlString)
	if err != nil {
		return nil, err
	}

	result := make([]StatWinRatio, 0)
	for rows.Next() {
		var swr StatWinRatio
		if err := rows.Scan(
			&swr.PlayCount,
			&swr.WinCount,
			&swr.WinRatio,
			&swr.Name); err != nil {
			log.Println(err)
			return result, errors.New("failed to scan row in query results")
		}
		result = append(result, swr)
	}

	return result, nil
}

func GetMostDrawsByPlayer() ([]StatCount, error) {
	sqlString := `
		SELECT
			COUNT(DISTINCT LD.ID) AS COUNT,
			U.NAME AS NAME
		FROM LOG_DRAW AS LD
			INNER JOIN USER AS U ON U.ID = LD.PLAYER_USER_ID
		GROUP BY U.ID
		ORDER BY
			COUNT DESC,
			NAME ASC
		LIMIT 5
	`
	rows, err := query(sqlString)
	if err != nil {
		return nil, err
	}

	result := make([]StatCount, 0)
	for rows.Next() {
		var sc StatCount
		if err := rows.Scan(&sc.Count, &sc.Name); err != nil {
			log.Println(err)
			return result, errors.New("failed to scan row in query results")
		}
		result = append(result, sc)
	}

	return result, nil
}

func GetMostDrawsByCard(userId uuid.UUID) ([]StatCount, error) {
	sqlString := `
		SELECT
			COUNT(DISTINCT LD.ID) AS COUNT,
			COALESCE(C.TEXT, LD.SPECIAL_CATEGORY, 'UNKNOWN') AS NAME
		FROM LOG_DRAW AS LD
			LEFT JOIN CARD AS C ON C.ID = LD.CARD_ID
		WHERE FN_USER_HAS_DECK_ACCESS(?, C.DECK_ID)
		GROUP BY C.ID
		ORDER BY
			COUNT DESC,
			NAME ASC
		LIMIT 5
	`
	rows, err := query(sqlString, userId)
	if err != nil {
		return nil, err
	}

	result := make([]StatCount, 0)
	for rows.Next() {
		var sc StatCount
		if err := rows.Scan(&sc.Count, &sc.Name); err != nil {
			log.Println(err)
			return result, errors.New("failed to scan row in query results")
		}
		result = append(result, sc)
	}

	return result, nil
}

func GetMostDrawsBySpecialCategory() ([]StatCount, error) {
	sqlString := `
		SELECT
			COUNT(DISTINCT LD.ID) AS COUNT,
			COALESCE(LD.SPECIAL_CATEGORY, 'NONE') AS NAME
		FROM LOG_DRAW AS LD
		GROUP BY LD.SPECIAL_CATEGORY
		ORDER BY
			COUNT DESC,
			NAME ASC
		LIMIT 5
	`
	rows, err := query(sqlString)
	if err != nil {
		return nil, err
	}

	result := make([]StatCount, 0)
	for rows.Next() {
		var sc StatCount
		if err := rows.Scan(&sc.Count, &sc.Name); err != nil {
			log.Println(err)
			return result, errors.New("failed to scan row in query results")
		}
		result = append(result, sc)
	}

	return result, nil
}

func GetMostDiscardsByPlayer() ([]StatCount, error) {
	sqlString := `
		SELECT
			COUNT(DISTINCT LD.ID) AS COUNT,
			U.NAME AS NAME
		FROM LOG_DISCARD AS LD
			INNER JOIN USER AS U ON U.ID = LD.PLAYER_USER_ID
		GROUP BY U.ID
		ORDER BY
			COUNT DESC,
			NAME ASC
		LIMIT 5
	`
	rows, err := query(sqlString)
	if err != nil {
		return nil, err
	}

	result := make([]StatCount, 0)
	for rows.Next() {
		var sc StatCount
		if err := rows.Scan(&sc.Count, &sc.Name); err != nil {
			log.Println(err)
			return result, errors.New("failed to scan row in query results")
		}
		result = append(result, sc)
	}

	return result, nil
}

func GetMostDiscardsByCard(userId uuid.UUID) ([]StatCount, error) {
	sqlString := `
		SELECT
			COUNT(DISTINCT LD.ID) AS COUNT,
			COALESCE(C.TEXT, 'UNKNOWN') AS NAME
		FROM LOG_DISCARD AS LD
			LEFT JOIN CARD AS C ON C.ID = LD.CARD_ID
		WHERE FN_USER_HAS_DECK_ACCESS(?, C.DECK_ID)
		GROUP BY C.ID
		ORDER BY
			COUNT DESC,
			NAME ASC
		LIMIT 5
	`
	rows, err := query(sqlString, userId)
	if err != nil {
		return nil, err
	}

	result := make([]StatCount, 0)
	for rows.Next() {
		var sc StatCount
		if err := rows.Scan(&sc.Count, &sc.Name); err != nil {
			log.Println(err)
			return result, errors.New("failed to scan row in query results")
		}
		result = append(result, sc)
	}

	return result, nil
}

func GetMostSkipsByPlayer() ([]StatCount, error) {
	sqlString := `
		SELECT
			COUNT(DISTINCT LS.ID) AS COUNT,
			U.NAME AS NAME
		FROM LOG_SKIP AS LS
			INNER JOIN USER AS U ON U.ID = LS.PLAYER_USER_ID
		GROUP BY U.ID
		ORDER BY
			COUNT DESC,
			NAME ASC
		LIMIT 5
	`
	rows, err := query(sqlString)
	if err != nil {
		return nil, err
	}

	result := make([]StatCount, 0)
	for rows.Next() {
		var sc StatCount
		if err := rows.Scan(&sc.Count, &sc.Name); err != nil {
			log.Println(err)
			return result, errors.New("failed to scan row in query results")
		}
		result = append(result, sc)
	}

	return result, nil
}

func GetMostSkipsByCard(userId uuid.UUID) ([]StatCount, error) {
	sqlString := `
		SELECT
			COUNT(DISTINCT LS.ID) AS COUNT,
			COALESCE(C.TEXT, 'UNKNOWN') AS NAME
		FROM LOG_SKIP AS LS
			LEFT JOIN CARD AS C ON C.ID = LS.CARD_ID
		WHERE FN_USER_HAS_DECK_ACCESS(?, C.DECK_ID)
		GROUP BY C.ID
		ORDER BY
			COUNT DESC,
			NAME ASC
		LIMIT 5
	`
	rows, err := query(sqlString, userId)
	if err != nil {
		return nil, err
	}

	result := make([]StatCount, 0)
	for rows.Next() {
		var sc StatCount
		if err := rows.Scan(&sc.Count, &sc.Name); err != nil {
			log.Println(err)
			return result, errors.New("failed to scan row in query results")
		}
		result = append(result, sc)
	}

	return result, nil
}
