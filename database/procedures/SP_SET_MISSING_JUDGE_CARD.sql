CREATE PROCEDURE IF NOT EXISTS SP_SET_MISSING_JUDGE_CARD(
    IN VAR_LOBBY_ID UUID
)
BEGIN
    IF (SELECT
            CARD_ID
        FROM JUDGE
        WHERE LOBBY_ID = VAR_LOBBY_ID) IS NULL THEN
        CALL SP_SET_NEXT_JUDGE_CARD(VAR_LOBBY_ID);
    END IF;
END;