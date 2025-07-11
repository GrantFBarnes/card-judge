CREATE PROCEDURE IF NOT EXISTS SP_SET_NEXT_JUDGE_PLAYER(IN VAR_LOBBY_ID UUID)
BEGIN
    DECLARE VAR_PLAYER_COUNT INT;
    DECLARE VAR_CURRENT_POSITION INT;
    DECLARE VAR_TRY_COUNT INT;
    DECLARE VAR_NEXT_POSITION INT;
    DECLARE VAR_NEXT_JUDGE_PLAYER_ID UUID;

    SELECT
        COUNT(*)
    INTO
        VAR_PLAYER_COUNT
    FROM PLAYER
    WHERE LOBBY_ID = VAR_LOBBY_ID;

    SELECT
        POSITION
    INTO
        VAR_CURRENT_POSITION
    FROM JUDGE
    WHERE LOBBY_ID = VAR_LOBBY_ID;

    SET VAR_TRY_COUNT = 0;
    SET VAR_NEXT_POSITION = VAR_CURRENT_POSITION;

    WHILE VAR_NEXT_JUDGE_PLAYER_ID IS NULL
        AND VAR_TRY_COUNT < VAR_PLAYER_COUNT
    DO
        SET VAR_TRY_COUNT = VAR_TRY_COUNT + 1;
        SET VAR_NEXT_POSITION = VAR_NEXT_POSITION + 1;

        IF VAR_NEXT_POSITION > VAR_PLAYER_COUNT THEN
            SET VAR_NEXT_POSITION = 1;
        END
        IF;

        SELECT
            ID
        INTO
            VAR_NEXT_JUDGE_PLAYER_ID
        FROM PLAYER
        WHERE IS_ACTIVE = 1
            AND LOBBY_ID = VAR_LOBBY_ID
            AND JOIN_ORDER = VAR_NEXT_POSITION;
    END
    WHILE;

    UPDATE JUDGE
    SET POSITION = VAR_NEXT_POSITION,
        PLAYER_ID = VAR_NEXT_JUDGE_PLAYER_ID,
        RESPONSE_COUNT = 1
    WHERE LOBBY_ID = VAR_LOBBY_ID;
END;