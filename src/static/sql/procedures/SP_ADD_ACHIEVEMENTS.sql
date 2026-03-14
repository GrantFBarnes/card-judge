CREATE
OR REPLACE PROCEDURE SP_ADD_ACHIEVEMENTS(IN VAR_USER_ID UUID)
BEGIN
    INSERT INTO USER_ACHIEVEMENT(
        USER_ID,
        ACHIEVEMENT_CATEGORY,
        CREDITS_SPENT_CATEGORY,
        ACHIEVED_AMOUNT
    )
    SELECT
        *
    FROM (
            SELECT
                VAR_USER_ID AS USER_ID,
                'WIN' AS ACHIEVEMENT_CATEGORY,
                '' AS CREDITS_SPENT_CATEGORY,
                COUNT(DISTINCT LOBBY_ID) AS ACHIEVED_AMOUNT
            FROM V_GAME_WINNER
            WHERE USER_ID = VAR_USER_ID
        ) AS T ON DUPLICATE KEY UPDATE ACHIEVED_AMOUNT = T.ACHIEVED_AMOUNT;

    INSERT INTO USER_ACHIEVEMENT(
        USER_ID,
        ACHIEVEMENT_CATEGORY,
        CREDITS_SPENT_CATEGORY,
        ACHIEVED_AMOUNT
    )
    SELECT
        *
    FROM (
            SELECT
                VAR_USER_ID AS USER_ID,
                'PLAY' AS ACHIEVEMENT_CATEGORY,
                '' AS CREDITS_SPENT_CATEGORY,
                COUNT(DISTINCT LOBBY_ID) AS ACHIEVED_AMOUNT
            FROM LOG_RESPONSE_CARD
            WHERE PLAYER_USER_ID = VAR_USER_ID
        ) AS T ON DUPLICATE KEY UPDATE ACHIEVED_AMOUNT = T.ACHIEVED_AMOUNT;

    INSERT INTO USER_ACHIEVEMENT(
        USER_ID,
        ACHIEVEMENT_CATEGORY,
        CREDITS_SPENT_CATEGORY,
        ACHIEVED_AMOUNT
    )
    SELECT
        *
    FROM (
            SELECT
                VAR_USER_ID AS USER_ID,
                'CREDITS-SPENT' AS ACHIEVEMENT_CATEGORY,
                CATEGORY AS CREDITS_SPENT_CATEGORY,
                COUNT(*) AS ACHIEVED_AMOUNT
            FROM LOG_CREDITS_SPENT
            WHERE USER_ID = VAR_USER_ID
                AND CATEGORY IN (
                    'WINNING-STREAK',
                    'LOSING-STREAK',
                    'SKIP-JUDGE',
                    'ALERT',
                    'GAMBLE',
                    'GAMBLE-WIN',
                    'BET',
                    'BET-WIN',
                    'EXTRA-RESPONSE',
                    'BLOCK-RESPONSE',
                    'STEAL',
                    'SURPRISE',
                    'FIND',
                    'WILD',
                    'PERK'
                )
            GROUP BY CATEGORY
        ) AS T ON DUPLICATE KEY UPDATE ACHIEVED_AMOUNT = T.ACHIEVED_AMOUNT;

    INSERT INTO USER_ACHIEVEMENT(
        USER_ID,
        ACHIEVEMENT_CATEGORY,
        CREDITS_SPENT_CATEGORY,
        ACHIEVED_AMOUNT
    )
    SELECT
        *
    FROM (
            SELECT
                VAR_USER_ID AS USER_ID,
                'DISCARD' AS ACHIEVEMENT_CATEGORY,
                '' AS CREDITS_SPENT_CATEGORY,
                COUNT(*) AS ACHIEVED_AMOUNT
            FROM LOG_DISCARD
            WHERE USER_ID = VAR_USER_ID
        ) AS T ON DUPLICATE KEY UPDATE ACHIEVED_AMOUNT = T.ACHIEVED_AMOUNT;

    INSERT INTO USER_ACHIEVEMENT(
        USER_ID,
        ACHIEVEMENT_CATEGORY,
        CREDITS_SPENT_CATEGORY,
        ACHIEVED_AMOUNT
    )
    SELECT
        *
    FROM (
            SELECT
                VAR_USER_ID AS USER_ID,
                'SKIP' AS ACHIEVEMENT_CATEGORY,
                '' AS CREDITS_SPENT_CATEGORY,
                COUNT(*) AS ACHIEVED_AMOUNT
            FROM LOG_SKIP
            WHERE USER_ID = VAR_USER_ID
        ) AS T ON DUPLICATE KEY UPDATE ACHIEVED_AMOUNT = T.ACHIEVED_AMOUNT;

    INSERT INTO USER_ACHIEVEMENT(
        USER_ID,
        ACHIEVEMENT_CATEGORY,
        CREDITS_SPENT_CATEGORY,
        ACHIEVED_AMOUNT
    )
    SELECT
        *
    FROM (
            SELECT
                VAR_USER_ID AS USER_ID,
                'KICK' AS ACHIEVEMENT_CATEGORY,
                '' AS CREDITS_SPENT_CATEGORY,
                COUNT(*) AS ACHIEVED_AMOUNT
            FROM LOG_KICK
            WHERE USER_ID = VAR_USER_ID
        ) AS T ON DUPLICATE KEY UPDATE ACHIEVED_AMOUNT = T.ACHIEVED_AMOUNT;

    INSERT INTO USER_ACHIEVEMENT(
        USER_ID,
        ACHIEVEMENT_CATEGORY,
        CREDITS_SPENT_CATEGORY,
        ACHIEVED_AMOUNT
    )
    SELECT
        *
    FROM (
            SELECT
                VAR_USER_ID AS USER_ID,
                'FLIP-TABLE' AS ACHIEVEMENT_CATEGORY,
                '' AS CREDITS_SPENT_CATEGORY,
                COUNT(*) AS ACHIEVED_AMOUNT
            FROM LOG_FLIP_TABLE
            WHERE USER_ID = VAR_USER_ID
        ) AS T ON DUPLICATE KEY UPDATE ACHIEVED_AMOUNT = T.ACHIEVED_AMOUNT;
END;