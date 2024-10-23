CREATE FUNCTION IF NOT EXISTS FN_GET_RANDOM_LOBBY_PLAYER (
    IN VAR_LOBBY_ID UUID
) RETURNS UUID
BEGIN
    RETURN (
        SELECT ID
        FROM PLAYER
        WHERE LOBBY_ID = VAR_LOBBY_ID
            AND IS_ACTIVE = 1
        ORDER BY RAND()
        LIMIT 1
    );
END;