CREATE PROCEDURE IF NOT EXISTS SP_WITHDRAW_LOBBY(
    IN VAR_LOBBY_ID UUID
)
BEGIN
    DECLARE VAR_LOOP_DONE BOOLEAN DEFAULT FALSE;
    DECLARE VAR_RESPONSE_ID UUID;
    DECLARE VAR_RESPONSE_CURSOR CURSOR FOR
        SELECT DISTINCT R.ID
        FROM RESPONSE AS R
                 INNER JOIN PLAYER AS P ON P.ID = R.PLAYER_ID
        WHERE P.LOBBY_ID = VAR_LOBBY_ID;
    DECLARE CONTINUE HANDLER FOR NOT FOUND SET VAR_LOOP_DONE = TRUE;

    OPEN VAR_RESPONSE_CURSOR;

    READ_LOOP:
    LOOP
        FETCH VAR_RESPONSE_CURSOR INTO VAR_RESPONSE_ID;
        IF VAR_LOOP_DONE THEN
            LEAVE READ_LOOP;
        END IF;

        CALL SP_WITHDRAW_RESPONSE(VAR_RESPONSE_ID);
    END LOOP;

    CLOSE VAR_RESPONSE_CURSOR;
END;