CREATE PROCEDURE IF NOT EXISTS SP_DISCARD_CARD(
    IN VAR_PLAYER_ID UUID,
    IN VAR_CARD_ID UUID
)
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

    IF EXISTS(SELECT
                  ID
              FROM HAND
              WHERE PLAYER_ID = VAR_PLAYER_ID
                AND CARD_ID = VAR_CARD_ID
                AND IS_LOCKED = 0) THEN
        DELETE
        FROM HAND
        WHERE PLAYER_ID = VAR_PLAYER_ID
          AND CARD_ID = VAR_CARD_ID;

        CALL SP_DRAW_HAND(VAR_PLAYER_ID);

        INSERT
        INTO LOG_DISCARD
            (
                LOBBY_ID,
                USER_ID,
                CARD_ID
            )
        VALUES
            (
                VAR_LOBBY_ID,
                VAR_PLAYER_USER_ID,
                VAR_CARD_ID
            );
    END IF;
END;