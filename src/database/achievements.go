package database

import (
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
)

type Achievement struct {
	Name     string
	Achieved bool
	Rarity   string
}

func GetAchievements(userId uuid.UUID) ([]Achievement, error) {
	var result []Achievement

	userCount, err := getUserCount()
	if err != nil {
		return result, err
	}

	sqlString := `
		WITH USER_ACHIEVEMENTS AS (
				SELECT
					ACHIEVEMENT_CODE
				FROM USER_ACHIEVEMENT
				WHERE USER_ID = ?
			)
		SELECT
			A.NAME,
			IF(UA.ACHIEVEMENT_CODE IS NOT NULL, 1, 0) AS ACHIEVED,
			(
				SELECT
					COUNT(*) / ?
				FROM USER_ACHIEVEMENT
				WHERE ACHIEVEMENT_CODE = A.CODE
			) AS RARITY
		FROM V_ACHIEVEMENT AS A
			LEFT JOIN USER_ACHIEVEMENTS AS UA ON UA.ACHIEVEMENT_CODE = A.CODE
		ORDER BY A.CODE
	`
	rows, err := query(sqlString, userId, userCount)
	if err != nil {
		return result, err
	}
	defer rows.Close()

	for rows.Next() {
		var achievement Achievement
		var rarity float32
		if err := rows.Scan(
			&achievement.Name,
			&achievement.Achieved,
			&rarity,
		); err != nil {
			log.Println(err)
			return result, errors.New("failed to scan row in query results")
		}

		achievement.Rarity = getPercentageString(rarity)

		result = append(result, achievement)
	}

	return result, nil
}

func getUserCount() (int, error) {
	var result int

	sqlString := `
		SELECT
			COUNT(*)
		FROM USER
	`
	rows, err := query(sqlString)
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

func getPercentageString(value float32) string {
	if value > 1 {
		value = 1.0
	}
	return fmt.Sprintf("%.1f%%", value*100)
}
