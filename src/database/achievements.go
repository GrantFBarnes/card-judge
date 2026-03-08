package database

import (
	"database/sql"
	"errors"
	"log"

	"github.com/google/uuid"
)

type Achievement struct {
	Name            string
	ProgressPercent float32
	DateAchieved    sql.NullTime
}

func GetAchievements(userId uuid.UUID) ([]Achievement, error) {
	var result []Achievement

	userWinRoundCount, err := getUserWinRoundCount(userId)
	if err != nil {
		return result, err
	}

	sqlString := `
		WITH USER_ACHIEVEMENTS AS (
				SELECT
					ACHIEVEMENTCODE,
					CREATED_ON_DATE
				FROM USER_ACHIEVEMENT
				WHERE USER_ID = ?
			)
		SELECT
			A.NAME,
			CASE
				WHEN A.CODE = 'WIN-ROUND-1' THEN ?/1
				WHEN A.CODE = 'WIN-ROUND-10' THEN ?/10
				WHEN A.CODE = 'WIN-ROUND-100' THEN ?/100
				WHEN A.CODE = 'WIN-ROUND-1000' THEN ?/1000
				ELSE 0
			END AS PROGRESS,
			UA.CREATED_ON_DATE
		FROM ACHIEVEMENT AS A
			LEFT JOIN USER_ACHIEVEMENTS AS UA ON UA.ACHIEVEMENTCODE = A.CODE
		ORDER BY UA.CREATED_ON_DATE DESC,
			A.CODE
	`
	rows, err := query(
		sqlString,
		userId,
		userWinRoundCount,
		userWinRoundCount,
		userWinRoundCount,
		userWinRoundCount,
	)
	if err != nil {
		return result, err
	}
	defer rows.Close()

	for rows.Next() {
		var achievement Achievement
		if err := rows.Scan(
			&achievement.Name,
			&achievement.ProgressPercent,
			&achievement.DateAchieved,
		); err != nil {
			log.Println(err)
			return result, errors.New("failed to scan row in query results")
		}

		achievement.ProgressPercent = achievement.ProgressPercent * 100
		if achievement.ProgressPercent > 100 {
			achievement.ProgressPercent = 100
		}

		result = append(result, achievement)
	}

	return result, nil
}

func getUserWinRoundCount(userId uuid.UUID) (int, error) {
	var result int

	sqlString := `
		SELECT
			COUNT(DISTINCT LRC.ROUND_ID)
		FROM LOG_RESPONSE_CARD AS LRC
			INNER JOIN LOG_WIN AS LW ON LW.RESPONSE_ID = LRC.RESPONSE_ID
		WHERE LRC.PLAYER_USER_ID = ?
	`
	rows, err := query(sqlString, userId)
	if err != nil {
		return result, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&result); err != nil {
			log.Println(err)
			return result, errors.New("failed to scan row in query results")
		}
	}

	return result, nil
}
