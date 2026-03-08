package database

import (
	"database/sql"
	"errors"
	"log"

	"github.com/google/uuid"
)

type Achievement struct {
	Name         string
	DateAchieved sql.NullTime
}

func GetAchievements(userId uuid.UUID) ([]Achievement, error) {
	var result []Achievement

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
			UA.CREATED_ON_DATE
		FROM ACHIEVEMENT AS A
			LEFT JOIN USER_ACHIEVEMENTS AS UA ON UA.ACHIEVEMENTCODE = A.CODE
		ORDER BY UA.CREATED_ON_DATE DESC,
			A.CODE
	`
	rows, err := query(sqlString, userId)
	if err != nil {
		return result, err
	}
	defer rows.Close()

	for rows.Next() {
		var achievement Achievement
		if err := rows.Scan(
			&achievement.Name,
			&achievement.DateAchieved,
		); err != nil {
			log.Println(err)
			return result, errors.New("failed to scan row in query results")
		}
		result = append(result, achievement)
	}

	return result, nil
}
