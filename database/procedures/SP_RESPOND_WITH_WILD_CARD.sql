CREATE PROCEDURE IF NOT EXISTS SP_RESPOND_WITH_WILD_CARD(
    IN VAR_PLAYER_ID UUID,
    IN VAR_CARD_TEXT VARCHAR(255)
)
BEGIN
    DECLARE VAR_LOBBY_ID UUID;
    DECLARE VAR_CARD_ID UUID;

    SET VAR_CARD_ID = UUID();

    SELECT
        LOBBY_ID
    INTO VAR_LOBBY_ID
    FROM PLAYER
    WHERE ID = VAR_PLAYER_ID;

    INSERT
    INTO CARD
        (
            ID,
            DECK_ID,
            CATEGORY,
            TEXT
        )
    SELECT
        VAR_CARD_ID   AS ID,
        D.ID          AS DECK_ID,
        'RESPONSE'    AS CATEGORY,
        VAR_CARD_TEXT AS TEXT
    FROM DECK AS D
    WHERE D.ID = VAR_LOBBY_ID
      AND D.IS_LOBBY_WILD_DECK = TRUE;

    INSERT
    INTO LOG_WILD
        (
            CARD_ID,
            CARD_TEXT
        )
    VALUES
        (
            VAR_CARD_ID,
            VAR_CARD_TEXT
        );

    UPDATE PLAYER
    SET
        CREDITS_SPENT = CREDITS_SPENT + 3
    WHERE ID = VAR_PLAYER_ID;

    INSERT
    INTO LOG_CREDITS_SPENT
        (
            LOBBY_ID,
            USER_ID,
            AMOUNT,
            CATEGORY
        )
    SELECT
        LOBBY_ID,
        USER_ID,
        3,
        'WILD'
    FROM PLAYER
    WHERE ID = VAR_PLAYER_ID;

    CALL SP_RESPOND_WITH_CARD(VAR_PLAYER_ID, VAR_CARD_ID, 'WILD');
END;