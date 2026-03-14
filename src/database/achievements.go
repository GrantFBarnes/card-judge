package database

import (
	"errors"
	"log"

	"github.com/google/uuid"
)

type Achievement struct {
	Name     string
	Progress int
}

func GetAchievements(userId uuid.UUID) ([]Achievement, error) {
	var result []Achievement

	sqlString := `
		SELECT
			CASE
				ACHIEVEMENT_CATEGORY
				WHEN 'WIN' THEN 'Games Won'
				WHEN 'PLAY' THEN 'Games Played'
				WHEN 'CREDITS-SPENT' THEN (
					CASE
						CREDITS_SPENT_CATEGORY
						WHEN 'WINNING-STREAK' THEN 'Winning Streaks'
						WHEN 'LOSING-STREAK' THEN 'Losing Streaks'
						WHEN 'SKIP-JUDGE' THEN 'Skipped Being Judge'
						WHEN 'ALERT' THEN 'Lobby Alerts'
						WHEN 'GAMBLE' THEN 'Gambles Made'
						WHEN 'GAMBLE-WIN' THEN 'Gambles Won'
						WHEN 'BET' THEN 'Bets Placed'
						WHEN 'BET-WIN' THEN 'Bets Won'
						WHEN 'EXTRA-RESPONSE' THEN 'Extra Responses'
						WHEN 'BLOCK-RESPONSE' THEN 'Blocked Responses'
						WHEN 'STEAL' THEN 'Stoken Cards Played'
						WHEN 'SURPRISE' THEN 'Surpise Cards Played'
						WHEN 'FIND' THEN 'Find Cards Played'
						WHEN 'WILD' THEN 'Wild Cards Played'
						WHEN 'PERK' THEN 'Perks Purchased'
						ELSE ''
					END
				)
				WHEN 'DISCARD' THEN 'Cards Discarded'
				WHEN 'SKIP' THEN 'Prompts Skipped'
				WHEN 'KICK' THEN 'Kicked From Lobby'
				WHEN 'FLIP-TABLE' THEN 'Flipped Tables'
				ELSE ''
			END AS ACHIEVEMENT_NAME,
			ACHIEVED_AMOUNT
		FROM USER_ACHIEVEMENT
		WHERE USER_ID = ?
		ORDER BY
			CASE
				ACHIEVEMENT_CATEGORY
				WHEN 'WIN' THEN 1
				WHEN 'PLAY' THEN 2
				WHEN 'CREDITS-SPENT' THEN (
					CASE
						CREDITS_SPENT_CATEGORY
						WHEN 'WINNING-STREAK' THEN 3
						WHEN 'LOSING-STREAK' THEN 4
						WHEN 'SKIP-JUDGE' THEN 5
						WHEN 'ALERT' THEN 6
						WHEN 'GAMBLE' THEN 7
						WHEN 'GAMBLE-WIN' THEN 8
						WHEN 'BET' THEN 9
						WHEN 'BET-WIN' THEN 10
						WHEN 'EXTRA-RESPONSE' THEN 11
						WHEN 'BLOCK-RESPONSE' THEN 12
						WHEN 'STEAL' THEN 13
						WHEN 'SURPRISE' THEN 14
						WHEN 'FIND' THEN 15
						WHEN 'WILD' THEN 16
						WHEN 'PERK' THEN 17
						ELSE 99
					END
				)
				WHEN 'DISCARD' THEN 18
				WHEN 'SKIP' THEN 19
				WHEN 'KICK' THEN 20
				WHEN 'FLIP-TABLE' THEN 21
				ELSE 99
			END
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
			&achievement.Progress,
		); err != nil {
			log.Println(err)
			return result, errors.New("failed to scan row in query results")
		}

		if achievement.Progress >= 100 {
			achievement.Progress = 100
		}

		result = append(result, achievement)
	}

	return result, nil
}
