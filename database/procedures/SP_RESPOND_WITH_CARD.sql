CREATE PROCEDURE IF NOT EXISTS SP_RESPOND_WITH_CARD(
    IN VAR_PLAYER_ID UUID,
    IN VAR_CARD_ID UUID,
    IN VAR_SPECIAL_CATEGORY ENUM ('STEAL','SURPRISE','WILD')
)
BEGIN
    DECLARE VAR_RESPONSE_ID UUID;
    DECLARE VAR_RESPONSE_CARD_ID UUID DEFAULT UUID();

    SELECT ID
    INTO VAR_RESPONSE_ID
    FROM RESPONSE
    WHERE PLAYER_ID = VAR_PLAYER_ID;

    IF VAR_RESPONSE_ID IS NULL THEN
        INSERT INTO RESPONSE (PLAYER_ID) VALUES (VAR_PLAYER_ID);
    END IF;

    SELECT ID
    INTO VAR_RESPONSE_ID
    FROM RESPONSE
    WHERE PLAYER_ID = VAR_PLAYER_ID;

    INSERT INTO RESPONSE_CARD (ID,
                               RESPONSE_ID,
                               CARD_ID,
                               SPECIAL_CATEGORY)
    VALUES (VAR_RESPONSE_CARD_ID,
            VAR_RESPONSE_ID,
            VAR_CARD_ID,
            VAR_SPECIAL_CATEGORY);

    INSERT INTO LOG_RESPONSE_CARD (LOBBY_ID,
                                   ROUND_ID,
                                   RESPONSE_ID,
                                   RESPONSE_CARD_ID,
                                   JUDGE_USER_ID,
                                   JUDGE_CARD_ID,
                                   PLAYER_USER_ID,
                                   PLAYER_CARD_ID,
                                   SPECIAL_CATEGORY)
    SELECT L.ID                AS LOBBY_ID,
           L.ROUND_ID          AS ROUND_ID,
           R.ID                AS RESPONSE_ID,
           RC.ID               AS RESPONSE_CARD_ID,
           J.PLAYER_ID         AS JUDGE_USER_ID,
           J.CARD_ID           AS JUDGE_CARD_ID,
           P.USER_ID           AS PLAYER_USER_ID,
           RC.CARD_ID          AS PLAYER_CARD_ID,
           RC.SPECIAL_CATEGORY AS SPECIAL_CATEGORY
    FROM RESPONSE_CARD AS RC
             INNER JOIN RESPONSE AS R ON R.ID = RC.RESPONSE_ID
             INNER JOIN PLAYER AS P ON P.ID = R.PLAYER_ID
             INNER JOIN LOBBY AS L ON L.ID = P.LOBBY_ID
             INNER JOIN JUDGE AS J ON J.LOBBY_ID = L.ID
    WHERE RC.ID = VAR_RESPONSE_CARD_ID;

    DELETE
    FROM HAND
    WHERE PLAYER_ID = VAR_PLAYER_ID
      AND CARD_ID = VAR_CARD_ID;
END;