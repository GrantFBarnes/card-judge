CREATE PROCEDURE IF NOT EXISTS SP_BET_ON_WIN(
    IN VAR_PLAYER_ID UUID,
    IN VAR_BET_ON_WIN INT
)
BEGIN
    UPDATE PLAYER
    SET
        CREDITS_SPENT = CREDITS_SPENT + VAR_BET_ON_WIN,
        BET_ON_WIN    = VAR_BET_ON_WIN
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
        VAR_BET_ON_WIN,
        'BET'
    FROM PLAYER
    WHERE ID = VAR_PLAYER_ID;
END;