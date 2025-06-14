CREATE PROCEDURE IF NOT EXISTS SP_SET_RESPONSE_COUNT (
    IN VAR_LOBBY_ID UUID,
    IN VAR_RESPONSE_COUNT INT
)
BEGIN
    UPDATE JUDGE
    SET RESPONSE_COUNT = VAR_RESPONSE_COUNT
    WHERE LOBBY_ID = VAR_LOBBY_ID;

    CALL SP_SET_RESPONSES_LOBBY(VAR_LOBBY_ID);
END;