CREATE PROCEDURE IF NOT EXISTS SP_SKIP_PROMPT (
    IN VAR_LOBBY_ID UUID
)
BEGIN
    DECLARE VAR_CARD_ID UUID;

    INSERT INTO LOG_SKIP (PLAYER_USER_ID, CARD_ID)
    SELECT
        P.USER_ID AS PLAYER_USER_ID,
        J.CARD_ID
    FROM JUDGE AS J
        INNER JOIN PLAYER AS P ON P.ID = J.PLAYER_ID
    WHERE J.LOBBY_ID = VAR_LOBBY_ID;

    SET VAR_CARD_ID = FN_GET_RANDOM_PROMPT_CARD (VAR_LOBBY_ID);

    IF VAR_CARD_ID IS NULL THEN
        DELETE FROM JUDGE
        WHERE LOBBY_ID = VAR_LOBBY_ID;
    ELSE
        UPDATE JUDGE
        SET CARD_ID = VAR_CARD_ID
        WHERE LOBBY_ID = VAR_LOBBY_ID;

        DELETE FROM DRAW_PILE
        WHERE CARD_ID = VAR_CARD_ID;
    END IF;
END;