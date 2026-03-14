package database

import (
	"errors"
	"log"

	"github.com/google/uuid"
)

type Achievement struct {
	Name     string
	Goal     int
	Achieved int
	Progress float32
	IsDone   bool
}

func GetAchievements(userId uuid.UUID) ([]Achievement, error) {
	var result []Achievement

	sqlString := `
		WITH USER_ACHIEVEMENTS AS (
				SELECT
					ACHIEVEMENT_CATEGORY,
					CREDITS_SPENT_CATEGORY,
					ACHIEVED_AMOUNT
				FROM USER_ACHIEVEMENT
				WHERE USER_ID = ?
			)
		SELECT
			A.ACHIEVEMENT_NAME,
			A.ACHIEVEMENT_AMOUNT,
			COALESCE(UA.ACHIEVED_AMOUNT, 0)
		FROM V_ACHIEVEMENT AS A
			LEFT JOIN USER_ACHIEVEMENTS AS UA ON UA.ACHIEVEMENT_CATEGORY = A.ACHIEVEMENT_CATEGORY
			AND UA.CREDITS_SPENT_CATEGORY = A.CREDITS_SPENT_CATEGORY
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
			&achievement.Goal,
			&achievement.Achieved,
		); err != nil {
			log.Println(err)
			return result, errors.New("failed to scan row in query results")
		}

		achievement.Progress = float32(achievement.Achieved) / float32(achievement.Goal) * 100.0
		if achievement.Progress >= 100.0 {
			achievement.Progress = 100.0
			achievement.IsDone = true
		}

		result = append(result, achievement)
	}

	return result, nil
}
