CREATE PROCEDURE IF NOT EXISTS SP_DISCARD_HAND (
    IN VAR_PLAYER_ID UUID
)
BEGIN
    DECLARE VAR_PLAYER_USER_ID UUID;

    SELECT USER_ID
    INTO VAR_PLAYER_USER_ID
    FROM PLAYER
    WHERE ID = VAR_PLAYER_ID;

    INSERT INTO LOG_DISCARD (PLAYER_USER_ID, CARD_ID)
    SELECT
        VAR_PLAYER_USER_ID AS PLAYER_USER_ID,
        CARD_ID
    FROM HAND
    WHERE PLAYER_ID = VAR_PLAYER_ID
        AND IS_LOCKED = 0;

    DELETE FROM HAND
    WHERE PLAYER_ID = VAR_PLAYER_ID
        AND IS_LOCKED = 0;
END;