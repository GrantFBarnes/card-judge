package database

import (
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
)

type Achievement struct {
	Name         string
	DateAchieved time.Time
}

func GetAchievements(userId uuid.UUID) ([]Achievement, error) {
	var result []Achievement

	sqlString := `
		SELECT
			UA.ACHIEVEMENT,
			UA.CREATED_ON_DATE
		FROM USER_ACHIEVEMENT AS UA
		WHERE UA.USER_ID = ?
		ORDER BY UA.CREATED_ON_DATE DESC
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
