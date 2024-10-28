CREATE TRIGGER IF NOT EXISTS TR_SET_RESPONSES_AF_UP_JUDGE
    AFTER UPDATE
    ON JUDGE
    FOR EACH ROW
BEGIN
    DECLARE VAR_PLAYER_RESPONSE_COUNT INT;
    DECLARE VAR_LOOP_DONE BOOLEAN DEFAULT FALSE;
    DECLARE VAR_PLAYER_ID UUID;
    DECLARE VAR_PLAYER_CURSOR CURSOR FOR
        SELECT ID
        FROM PLAYER
        WHERE LOBBY_ID = NEW.LOBBY_ID;
    DECLARE CONTINUE HANDLER FOR NOT FOUND SET VAR_LOOP_DONE = TRUE;

    IF NEW.RESPONSE_COUNT > OLD.RESPONSE_COUNT THEN
        BEGIN
            OPEN VAR_PLAYER_CURSOR;

            READ_LOOP:
            LOOP
                FETCH VAR_PLAYER_CURSOR INTO VAR_PLAYER_ID;
                IF VAR_LOOP_DONE THEN
                    LEAVE READ_LOOP;
                END IF;

                SELECT COUNT(*)
                INTO VAR_PLAYER_RESPONSE_COUNT
                FROM RESPONSE
                WHERE PLAYER_ID = VAR_PLAYER_ID;

                -- CREATE ANY MISSING RESPONSES
                WHILE VAR_PLAYER_RESPONSE_COUNT < NEW.RESPONSE_COUNT
                    DO
                        INSERT INTO RESPONSE (PLAYER_ID) VALUE (VAR_PLAYER_ID);

                        SELECT COUNT(*)
                        INTO VAR_PLAYER_RESPONSE_COUNT
                        FROM RESPONSE
                        WHERE PLAYER_ID = VAR_PLAYER_ID;
                    END WHILE;
            END LOOP;

            CLOSE VAR_PLAYER_CURSOR;
        END;
    ELSEIF NEW.RESPONSE_COUNT < OLD.RESPONSE_COUNT THEN
        BEGIN
            OPEN VAR_PLAYER_CURSOR;

            READ_LOOP:
            LOOP
                FETCH VAR_PLAYER_CURSOR INTO VAR_PLAYER_ID;
                IF VAR_LOOP_DONE THEN
                    LEAVE READ_LOOP;
                END IF;

                SELECT COUNT(*)
                INTO VAR_PLAYER_RESPONSE_COUNT
                FROM RESPONSE
                WHERE PLAYER_ID = VAR_PLAYER_ID;

                -- DELETE ANY EXTRA RESPONSES
                WHILE VAR_PLAYER_RESPONSE_COUNT > NEW.RESPONSE_COUNT
                    DO
                        DELETE
                        FROM RESPONSE
                        WHERE ID IN (SELECT ID
                                     FROM (SELECT ID
                                           FROM RESPONSE
                                           WHERE PLAYER_ID = VAR_PLAYER_ID
                                           ORDER BY CREATED_ON_DATE DESC
                                           LIMIT 1) AS T);

                        SELECT COUNT(*)
                        INTO VAR_PLAYER_RESPONSE_COUNT
                        FROM RESPONSE
                        WHERE PLAYER_ID = VAR_PLAYER_ID;
                    END WHILE;
            END LOOP;

            CLOSE VAR_PLAYER_CURSOR;
        END;
    END IF;
END;
