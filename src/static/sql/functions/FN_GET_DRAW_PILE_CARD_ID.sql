CREATE
OR REPLACE FUNCTION FN_GET_DRAW_PILE_CARD_ID(
    IN VAR_CATEGORY ENUM('PROMPT', 'RESPONSE'),
    IN VAR_LOBBY_ID UUID
)
RETURNS UUID
BEGIN
    DECLARE VAR_LOBBY_DRAW_PRIORITY ENUM('RANDOM', 'PLAYCOUNT') DEFAULT (
            SELECT
                DRAW_PRIORITY
            FROM LOBBY
            WHERE ID = VAR_LOBBY_ID
        );

    IF VAR_LOBBY_DRAW_PRIORITY = 'PLAYCOUNT' THEN
        IF VAR_CATEGORY = 'PROMPT' THEN
            RETURN (
                SELECT
                    ID
                FROM (
                        SELECT
                            C.ID,
                            (
                                SELECT
                                    COUNT(DISTINCT ROUND_ID)
                                FROM LOG_RESPONSE_CARD
                                WHERE JUDGE_CARD_ID = C.ID
                            ) AS PLAY_COUNT,
                            (SELECT COUNT(*) FROM LOG_SKIP WHERE CARD_ID = C.ID) AS SKIP_COUNT
                        FROM DRAW_PILE AS DP
                            INNER JOIN CARD AS C ON C.ID = DP.CARD_ID
                        WHERE C.CATEGORY = VAR_CATEGORY
                            AND DP.LOBBY_ID = VAR_LOBBY_ID
                    ) AS T
                ORDER BY PLAY_COUNT,
                    SKIP_COUNT,
                    RAND()
                LIMIT 1
            );
        END IF;

        RETURN (
            SELECT
                ID
            FROM (
                    SELECT
                        C.ID,
                        (
                            SELECT
                                COUNT(*)
                            FROM LOG_RESPONSE_CARD
                            WHERE PLAYER_CARD_ID = C.ID
                        ) AS PLAY_COUNT,
                        (SELECT COUNT(*) FROM LOG_DISCARD WHERE CARD_ID = C.ID) AS DISCARD_COUNT
                    FROM DRAW_PILE AS DP
                        INNER JOIN CARD AS C ON C.ID = DP.CARD_ID
                    WHERE C.CATEGORY = VAR_CATEGORY
                        AND DP.LOBBY_ID = VAR_LOBBY_ID
                ) AS T
            ORDER BY PLAY_COUNT,
                DISCARD_COUNT,
                RAND()
            LIMIT 1
        );
    END IF;

    RETURN (
        SELECT
            C.ID
        FROM DRAW_PILE AS DP
            INNER JOIN CARD AS C ON C.ID = DP.CARD_ID
        WHERE C.CATEGORY = VAR_CATEGORY
            AND DP.LOBBY_ID = VAR_LOBBY_ID
        ORDER BY RAND()
        LIMIT 1
    );
END;
