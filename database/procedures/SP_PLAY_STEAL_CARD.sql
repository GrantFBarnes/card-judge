CREATE PROCEDURE IF NOT EXISTS SP_PLAY_STEAL_CARD (
    IN VAR_PLAYER_ID UUID
)
BEGIN
    DECLARE VAR_LOBBY_ID UUID;
    DECLARE VAR_PLAYER_USER_ID UUID;
    DECLARE VAR_VICTIM_PLAYER_ID UUID;
    DECLARE VAR_CARD_ID UUID;

    SELECT
        LOBBY_ID,
        USER_ID
    INTO
        VAR_LOBBY_ID,
        VAR_PLAYER_USER_ID
    FROM PLAYER
    WHERE ID = VAR_PLAYER_ID;

    SELECT
        P.ID,
        C.ID
    INTO
        VAR_VICTIM_PLAYER_ID,
        VAR_CARD_ID
    FROM PLAYER AS P
        INNER JOIN HAND AS H ON H.PLAYER_ID = P.ID
        INNER JOIN CARD AS C ON C.ID = H.CARD_ID
    WHERE P.LOBBY_ID = VAR_LOBBY_ID
        AND P.ID <> VAR_PLAYER_ID
    ORDER BY RAND()
    LIMIT 1;

    UPDATE PLAYER
    SET CREDITS_SPENT = CREDITS_SPENT + 1
    WHERE ID = VAR_PLAYER_ID;

    DELETE FROM HAND
    WHERE PLAYER_ID = VAR_VICTIM_PLAYER_ID
        AND CARD_ID = VAR_CARD_ID;

    UPDATE PLAYER
    SET CREDITS_SPENT = CREDITS_SPENT - 1
    WHERE ID = VAR_VICTIM_PLAYER_ID;

    INSERT INTO LOG_DRAW (PLAYER_USER_ID, CARD_ID, SPECIAL_CATEGORY)
    VALUES (VAR_PLAYER_USER_ID, VAR_CARD_ID, 'STEAL');

    CALL SP_PLAY_CARD (VAR_PLAYER_ID, VAR_CARD_ID, 'STEAL');
END;