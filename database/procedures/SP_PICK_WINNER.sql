CREATE PROCEDURE IF NOT EXISTS SP_PICK_WINNER(
    IN VAR_RESPONSE_ID UUID
)
BEGIN
    DECLARE VAR_PLAYER_ID UUID;
    DECLARE VAR_PLAYER_BET_ON_WIN INT;
    DECLARE VAR_LOBBY_ID UUID;

    SELECT
        P.ID AS PLAYER_ID,
        P.BET_ON_WIN,
        P.LOBBY_ID
    INTO
        VAR_PLAYER_ID,
        VAR_PLAYER_BET_ON_WIN,
        VAR_LOBBY_ID
    FROM RESPONSE AS R
             INNER JOIN PLAYER AS P ON P.ID = R.PLAYER_ID
    WHERE R.ID = VAR_RESPONSE_ID;

    IF (VAR_PLAYER_BET_ON_WIN > 0) THEN
        UPDATE PLAYER
        SET
            CREDITS_SPENT = CREDITS_SPENT - (VAR_PLAYER_BET_ON_WIN * 2)
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
            (VAR_PLAYER_BET_ON_WIN * 2) * -1,
            'BET-WIN'
        FROM PLAYER
        WHERE ID = VAR_PLAYER_ID;
    END IF;

    INSERT
    INTO WIN
        (
            PLAYER_ID
        )
    VALUES
        (
            VAR_PLAYER_ID
        );
    INSERT
    INTO LOG_WIN
        (
            RESPONSE_ID
        )
    VALUES
        (
            VAR_RESPONSE_ID
        );

    CALL SP_SET_WINNING_STREAK(VAR_PLAYER_ID);
    CALL SP_SET_LOSING_STREAK(VAR_PLAYER_ID);
    CALL SP_START_NEW_ROUND(VAR_LOBBY_ID);

    SELECT
        U.NAME
    FROM PLAYER AS P
             INNER JOIN USER AS U ON U.ID = P.USER_ID
    WHERE P.ID = VAR_PLAYER_ID;
END;