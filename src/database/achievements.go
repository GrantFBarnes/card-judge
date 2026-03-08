package database

import (
	"errors"
	"log"

	"github.com/google/uuid"
)

type Achievement struct {
	Category   string
	GoalAmount int
	IsDone     bool
	Rarity     float32
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
			A.CATEGORY,
			A.GOAL_AMOUNT,
			IF(UA.ACHIEVEMENT_CODE IS NOT NULL, 1, 0) AS IS_DONE,
			(SELECT COUNT(*) FROM USER_ACHIEVEMENT WHERE ACHIEVEMENT_CODE = A.CODE) AS DONE_COUNT
		FROM V_ACHIEVEMENT AS A
			LEFT JOIN USER_ACHIEVEMENTS AS UA ON UA.ACHIEVEMENT_CODE = A.CODE
		ORDER BY A.LIST_ORDER
	`
	rows, err := query(sqlString, userId)
	if err != nil {
		return result, err
	}
	defer rows.Close()

	for rows.Next() {
		var achievement Achievement
		var doneCount int
		if err := rows.Scan(
			&achievement.Category,
			&achievement.GoalAmount,
			&achievement.IsDone,
			&doneCount,
		); err != nil {
			log.Println(err)
			return result, errors.New("failed to scan row in query results")
		}

		achievement.Rarity = float32(doneCount) / float32(userCount) * 100.0

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
