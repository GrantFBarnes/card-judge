CREATE PROCEDURE IF NOT EXISTS SP_DRAW_HAND (
    IN VAR_PLAYER_ID UUID
)
BEGIN
    DECLARE VAR_LOBBY_HAND_SIZE INT;
    DECLARE VAR_PLAYER_HAND_SIZE INT;
    DECLARE VAR_CARDS_TO_DRAW INT;

    SELECT L.HAND_SIZE
    INTO VAR_LOBBY_HAND_SIZE
    FROM LOBBY AS L
        INNER JOIN PLAYER AS P ON P.LOBBY_ID = L.ID
    WHERE P.ID = VAR_PLAYER_ID;

    SELECT COUNT(CARD_ID)
    INTO VAR_PLAYER_HAND_SIZE
    FROM HAND
    WHERE PLAYER_ID = VAR_PLAYER_ID;

    SET VAR_CARDS_TO_DRAW = VAR_LOBBY_HAND_SIZE - VAR_PLAYER_HAND_SIZE;

    INSERT INTO HAND (
        PLAYER_ID,
        CARD_ID
    )
    SELECT
        P.ID AS PLAYER_ID,
        C.ID AS CARD_ID
    FROM DRAW_PILE AS DP
        INNER JOIN PLAYER AS P ON P.LOBBY_ID = DP.LOBBY_ID
        INNER JOIN CARD AS C ON C.ID = DP.CARD_ID
    WHERE C.CATEGORY = 'RESPONSE'
        AND P.ID = VAR_PLAYER_ID
    ORDER BY RAND()
    LIMIT VAR_CARDS_TO_DRAW;

    DELETE DP
    FROM DRAW_PILE AS DP
        INNER JOIN PLAYER AS P ON P.LOBBY_ID = DP.LOBBY_ID
        INNER JOIN HAND AS H ON H.PLAYER_ID = P.ID AND H.CARD_ID = DP.CARD_ID
    WHERE P.ID = VAR_PLAYER_ID;

    INSERT INTO LOG_DRAW (
        PLAYER_USER_ID,
        CARD_ID
    )
    SELECT
        P.USER_ID AS PLAYER_USER_ID,
        H.CARD_ID
    FROM HAND AS H
        INNER JOIN PLAYER AS P ON P.ID = H.PLAYER_ID
    WHERE P.ID = VAR_PLAYER_ID
    ORDER BY H.CREATED_ON_DATE DESC
    LIMIT VAR_CARDS_TO_DRAW;
END;