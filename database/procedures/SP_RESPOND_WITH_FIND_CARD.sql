CREATE PROCEDURE IF NOT EXISTS SP_RESPOND_WITH_FIND_CARD(
    IN VAR_PLAYER_ID UUID,
    IN VAR_CARD_ID UUID
)
BEGIN
    DECLARE VAR_LOBBY_ID UUID;

    SELECT
        LOBBY_ID
    INTO
        VAR_LOBBY_ID
    FROM PLAYER
    WHERE ID = VAR_PLAYER_ID;

    IF EXISTS(
        SELECT
            ID
        FROM DRAW_PILE
        WHERE LOBBY_ID = VAR_LOBBY_ID
            AND CARD_ID = VAR_CARD_ID
    ) THEN
        DELETE
        FROM DRAW_PILE
        WHERE LOBBY_ID = VAR_LOBBY_ID
            AND CARD_ID = VAR_CARD_ID;

        UPDATE PLAYER
        SET CREDITS_SPENT = CREDITS_SPENT + 2
        WHERE ID = VAR_PLAYER_ID;

        INSERT INTO LOG_CREDITS_SPENT(LOBBY_ID, USER_ID, AMOUNT, CATEGORY)
        SELECT
            LOBBY_ID,
            USER_ID,
            2,
            'FIND'
        FROM PLAYER
        WHERE ID = VAR_PLAYER_ID;

        CALL SP_RESPOND_WITH_CARD(VAR_PLAYER_ID, VAR_CARD_ID, 'FIND');
    END
    IF;
END;