CREATE PROCEDURE IF NOT EXISTS SP_GAMBLE_CREDITS(
    IN VAR_PLAYER_ID UUID,
    IN VAR_CREDIT_AMOUNT INT
)
BEGIN
    UPDATE PLAYER
    SET CREDITS_SPENT = CREDITS_SPENT + VAR_CREDIT_AMOUNT
    WHERE ID = VAR_PLAYER_ID;

    INSERT INTO LOG_CREDITS_SPENT(LOBBY_ID, USER_ID, AMOUNT, CATEGORY)
    SELECT
        LOBBY_ID,
        USER_ID,
        VAR_CREDIT_AMOUNT,
        'GAMBLE'
    FROM PLAYER
    WHERE ID = VAR_PLAYER_ID;

    IF RAND() > 0.5 THEN
        UPDATE PLAYER
        SET CREDITS_SPENT = CREDITS_SPENT - (VAR_CREDIT_AMOUNT * 2)
        WHERE ID = VAR_PLAYER_ID;

        INSERT INTO LOG_CREDITS_SPENT(LOBBY_ID, USER_ID, AMOUNT, CATEGORY)
        SELECT
            LOBBY_ID,
            USER_ID,
            (VAR_CREDIT_AMOUNT * 2) * -1,
            'GAMBLE-WIN'
        FROM PLAYER
        WHERE ID = VAR_PLAYER_ID;
    END
    IF;
END;