CREATE FUNCTION IF NOT EXISTS FN_GET_PLAYER_IS_WINNING(IN VAR_PLAYER_ID UUID)
RETURNS BOOLEAN
BEGIN
    DECLARE VAR_LOBBY_ID UUID;
    DECLARE VAR_PLAYER_USER_ID UUID;

    SELECT
        LOBBY_ID,
        USER_ID
    INTO
        VAR_LOBBY_ID,
        VAR_PLAYER_USER_ID
    FROM PLAYER
    WHERE ID = VAR_PLAYER_ID;

    RETURN EXISTS(
        SELECT
            1
        FROM (
                SELECT
                    LRC.LOBBY_ID,
                    LRC.PLAYER_USER_ID,
                    RANK() OVER (
                        PARTITION BY LRC.LOBBY_ID
                        ORDER BY COUNT(DISTINCT LRC.ROUND_ID) DESC
                    ) AS RANKING
                FROM LOG_RESPONSE_CARD AS LRC
                    INNER JOIN LOG_WIN AS LW ON LW.RESPONSE_ID = LRC.RESPONSE_ID
                GROUP BY LRC.LOBBY_ID,
                    LRC.PLAYER_USER_ID
            ) AS ROUND_WINS
        WHERE LOBBY_ID = VAR_LOBBY_ID
            AND PLAYER_USER_ID = VAR_PLAYER_USER_ID
            AND RANKING = 1
    );
END;