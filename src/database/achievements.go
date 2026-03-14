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
		WITH CREDITS_SPENT AS (
				SELECT
					CATEGORY
				FROM LOG_CREDITS_SPENT
				WHERE USER_ID = ?
			)
		SELECT
			'Games Won',
			(SELECT COUNT(DISTINCT LOBBY_ID) FROM V_GAME_WINNER WHERE USER_ID = ?)
		UNION
		SELECT
			'Games Played',
			(SELECT COUNT(DISTINCT LOBBY_ID) FROM LOG_RESPONSE_CARD WHERE PLAYER_USER_ID = ?)
		UNION
		SELECT
			'Winning Streaks',
			(SELECT COUNT(*) FROM CREDITS_SPENT WHERE CATEGORY = 'WINNING-STREAK')
		UNION
		SELECT
			'Losing Streaks',
			(SELECT COUNT(*) FROM CREDITS_SPENT WHERE CATEGORY = 'LOSING-STREAK')
		UNION
		SELECT
			'Skipped Being Judge',
			(SELECT COUNT(*) FROM CREDITS_SPENT WHERE CATEGORY = 'SKIP-JUDGE')
		UNION
		SELECT
			'Lobby Alerts',
			(SELECT COUNT(*) FROM CREDITS_SPENT WHERE CATEGORY = 'ALERT')
		UNION
		SELECT
			'Gambles Made',
			(SELECT COUNT(*) FROM CREDITS_SPENT WHERE CATEGORY = 'GAMBLE')
		UNION
		SELECT
			'Gambles Won',
			(SELECT COUNT(*) FROM CREDITS_SPENT WHERE CATEGORY = 'GAMBLE-WIN')
		UNION
		SELECT
			'Bets Placed',
			(SELECT COUNT(*) FROM CREDITS_SPENT WHERE CATEGORY = 'BET')
		UNION
		SELECT
			'Bets Won',
			(SELECT COUNT(*) FROM CREDITS_SPENT WHERE CATEGORY = 'BET-WIN')
		UNION
		SELECT
			'Extra Responses',
			(SELECT COUNT(*) FROM CREDITS_SPENT WHERE CATEGORY = 'EXTRA-RESPONSE')
		UNION
		SELECT
			'Blocked Responses',
			(SELECT COUNT(*) FROM CREDITS_SPENT WHERE CATEGORY = 'BLOCK-RESPONSE')
		UNION
		SELECT
			'Stoken Cards Played',
			(SELECT COUNT(*) FROM CREDITS_SPENT WHERE CATEGORY = 'STEAL')
		UNION
		SELECT
			'Surpise Cards Played',
			(SELECT COUNT(*) FROM CREDITS_SPENT WHERE CATEGORY = 'SURPRISE')
		UNION
		SELECT
			'Find Cards Played',
			(SELECT COUNT(*) FROM CREDITS_SPENT WHERE CATEGORY = 'FIND')
		UNION
		SELECT
			'Wild Cards Played',
			(SELECT COUNT(*) FROM CREDITS_SPENT WHERE CATEGORY = 'WILD')
		UNION
		SELECT
			'Perks Purchased',
			(SELECT COUNT(*) FROM CREDITS_SPENT WHERE CATEGORY = 'PERK')
		UNION
		SELECT
			'Cards Discarded',
			(SELECT COUNT(*) FROM LOG_DISCARD WHERE USER_ID = ?)
		UNION
		SELECT
			'Prompts Skipped',
			(SELECT COUNT(*) FROM LOG_SKIP WHERE USER_ID = ?)
		UNION
		SELECT
			'Kicked From Lobby',
			(SELECT COUNT(*) FROM LOG_KICK WHERE USER_ID = ?)
		UNION
		SELECT
			'Flipped Tables',
			(SELECT COUNT(*) FROM LOG_FLIP_TABLE WHERE USER_ID = ?)
	`
	rows, err := query(sqlString, userId, userId, userId, userId, userId, userId, userId)
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
