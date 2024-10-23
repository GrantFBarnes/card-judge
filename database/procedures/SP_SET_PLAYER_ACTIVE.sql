CREATE PROCEDURE IF NOT EXISTS SP_SET_PLAYER_ACTIVE (
    IN VAR_PLAYER_ID UUID,
    IN VAR_LOBBY_ID UUID,
    IN VAR_USER_ID UUID
)
BEGIN
    IF EXISTS(
        SELECT ID
        FROM PLAYER
        WHERE ID = VAR_PLAYER_ID
    ) THEN
        UPDATE PLAYER
        SET IS_ACTIVE = 1
        WHERE ID = VAR_PLAYER_ID;
    ELSE
        INSERT INTO PLAYER (ID, LOBBY_ID, USER_ID)
        VALUES (VAR_PLAYER_ID, VAR_LOBBY_ID, VAR_USER_ID);
    END IF;
END;