CREATE PROCEDURE IF NOT EXISTS SP_PLAY_WILD_CARD (
    IN VAR_PLAYER_ID UUID,
    IN VAR_CARD_TEXT VARCHAR(255)
)
BEGIN
    DECLARE VAR_LOBBY_ID UUID;
    DECLARE VAR_PLAYER_USER_ID UUID;
    DECLARE VAR_CARD_ID UUID;

    SET VAR_CARD_ID = UUID();

    SELECT
        LOBBY_ID,
        USER_ID
    INTO
        VAR_LOBBY_ID,
        VAR_PLAYER_USER_ID
    FROM PLAYER
    WHERE ID = VAR_PLAYER_ID;

    INSERT INTO CARD (ID, DECK_ID, CATEGORY, TEXT)
    VALUES (VAR_CARD_ID, '00000000-0000-0000-0000-000000000000', 'RESPONSE', VAR_CARD_TEXT);

    UPDATE PLAYER
    SET CREDITS_SPENT = CREDITS_SPENT + 2
    WHERE ID = VAR_PLAYER_ID;

    INSERT INTO LOG_DRAW (LOBBY_ID, PLAYER_USER_ID, CARD_ID, SPECIAL_CATEGORY)
    VALUES (VAR_LOBBY_ID, VAR_PLAYER_USER_ID, VAR_CARD_ID, 'WILD');

    CALL SP_PLAY_CARD (VAR_PLAYER_ID, VAR_CARD_ID, 'WILD');
END;