CREATE PROCEDURE IF NOT EXISTS SP_SET_RESPONSES_PLAYER(
    IN VAR_PLAYER_ID UUID
)
BEGIN
    DECLARE VAR_LOBBY_ID UUID;
    DECLARE VAR_EXTRA_RESPONSES INT;
    DECLARE VAR_RESPONSE_COUNT_FINAL INT;
    DECLARE VAR_RESPONSE_COUNT_CURRENT INT;
    DECLARE VAR_RESPONSE_ID_TO_REMOVE UUID;

    SELECT
        LOBBY_ID,
        EXTRA_RESPONSES
    INTO
        VAR_LOBBY_ID,
        VAR_EXTRA_RESPONSES
    FROM PLAYER
    WHERE ID = VAR_PLAYER_ID;

    -- GET STANDARD RESPONSE COUNT FINAL
    SELECT
        RESPONSE_COUNT
    INTO VAR_RESPONSE_COUNT_FINAL
    FROM JUDGE
    WHERE LOBBY_ID = VAR_LOBBY_ID;

    -- ADD ANY EXTRA RESPONSES
    SET VAR_RESPONSE_COUNT_FINAL = VAR_RESPONSE_COUNT_FINAL + VAR_EXTRA_RESPONSES;

    -- CLEAR RESPONSES FOR INACTIVE PLAYER
    IF EXISTS(SELECT ID FROM PLAYER WHERE ID = VAR_PLAYER_ID AND IS_ACTIVE = 0) THEN
        SET VAR_RESPONSE_COUNT_FINAL = 0;
    END IF;

    -- CLEAR RESPONSES FOR JUDGE
    IF EXISTS(SELECT ID FROM JUDGE WHERE PLAYER_ID = VAR_PLAYER_ID) THEN
        SET VAR_RESPONSE_COUNT_FINAL = 0;
    END IF;

    SELECT
        COUNT(*)
    INTO VAR_RESPONSE_COUNT_CURRENT
    FROM RESPONSE
    WHERE PLAYER_ID = VAR_PLAYER_ID;

    -- CREATE ANY MISSING RESPONSES
    WHILE VAR_RESPONSE_COUNT_CURRENT < VAR_RESPONSE_COUNT_FINAL
        DO
            INSERT
            INTO RESPONSE
                (
                    PLAYER_ID
                )
            VALUES
                (
                    VAR_PLAYER_ID
                );

            SELECT
                COUNT(*)
            INTO VAR_RESPONSE_COUNT_CURRENT
            FROM RESPONSE
            WHERE PLAYER_ID = VAR_PLAYER_ID;
        END WHILE;

    -- DELETE ANY EXTRA RESPONSES
    WHILE VAR_RESPONSE_COUNT_CURRENT > VAR_RESPONSE_COUNT_FINAL
        DO
            SELECT
                ID
            INTO VAR_RESPONSE_ID_TO_REMOVE
            FROM RESPONSE
            WHERE PLAYER_ID = VAR_PLAYER_ID
            ORDER BY CREATED_ON_DATE DESC
            LIMIT 1;

            CALL SP_WITHDRAW_RESPONSE(VAR_RESPONSE_ID_TO_REMOVE);

            DELETE
            FROM RESPONSE
            WHERE ID = VAR_RESPONSE_ID_TO_REMOVE;

            SELECT
                COUNT(*)
            INTO VAR_RESPONSE_COUNT_CURRENT
            FROM RESPONSE
            WHERE PLAYER_ID = VAR_PLAYER_ID;
        END WHILE;
END;