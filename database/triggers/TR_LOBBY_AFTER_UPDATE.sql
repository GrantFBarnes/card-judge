CREATE TRIGGER IF NOT EXISTS TR_LOBBY_AFTER_UPDATE
AFTER UPDATE ON LOBBY
FOR EACH ROW
BEGIN
    IF NEW.HAND_SIZE > OLD.HAND_SIZE THEN
        BEGIN
            DECLARE VAR_LOOP_DONE BOOLEAN DEFAULT FALSE;
            DECLARE VAR_PLAYER_ID UUID;

            DECLARE VAR_PLAYER_CURSOR CURSOR
            FOR
            SELECT
                ID
            FROM PLAYER
            WHERE LOBBY_ID = NEW.ID;

            DECLARE CONTINUE HANDLER
            FOR NOT FOUND
            SET VAR_LOOP_DONE = TRUE;

            OPEN VAR_PLAYER_CURSOR;

                READ_LOOP: LOOP
                FETCH VAR_PLAYER_CURSOR
                INTO
                    VAR_PLAYER_ID;

                IF VAR_LOOP_DONE THEN LEAVE READ_LOOP;
                END
                IF;

                CALL SP_DRAW_HAND(VAR_PLAYER_ID);
                END LOOP;
            CLOSE VAR_PLAYER_CURSOR;
        END;
    END
    IF;
END;