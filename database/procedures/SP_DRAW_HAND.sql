CREATE PROCEDURE IF NOT EXISTS SP_DRAW_HAND(
    IN VAR_PLAYER_ID UUID
)
BEGIN
    DECLARE VAR_LOBBY_ID UUID;
    DECLARE VAR_PLAYER_USER_ID UUID;
    DECLARE VAR_LOBBY_HAND_SIZE INT;
    DECLARE VAR_LOBBY_DRAW_PILE_SIZE INT;
    DECLARE VAR_PLAYER_HAND_SIZE INT;
    DECLARE VAR_CARD_ID UUID;

    SELECT
        LOBBY_ID,
        USER_ID
    INTO
        VAR_LOBBY_ID,
        VAR_PLAYER_USER_ID
    FROM PLAYER
    WHERE ID = VAR_PLAYER_ID;

    SELECT
        HAND_SIZE
    INTO VAR_LOBBY_HAND_SIZE
    FROM LOBBY
    WHERE ID = VAR_LOBBY_ID;

    SELECT
        COUNT(C.ID)
    INTO VAR_LOBBY_DRAW_PILE_SIZE
    FROM DRAW_PILE AS DP
             INNER JOIN CARD AS C ON C.ID = DP.CARD_ID
    WHERE C.CATEGORY = 'RESPONSE'
      AND DP.LOBBY_ID = VAR_LOBBY_ID;

    SELECT
        COUNT(CARD_ID)
    INTO VAR_PLAYER_HAND_SIZE
    FROM HAND
    WHERE PLAYER_ID = VAR_PLAYER_ID;

    WHILE VAR_LOBBY_DRAW_PILE_SIZE > 0 AND VAR_LOBBY_HAND_SIZE > VAR_PLAYER_HAND_SIZE
        DO
            SELECT
                C.ID
            INTO VAR_CARD_ID
            FROM DRAW_PILE AS DP
                     INNER JOIN CARD AS C ON C.ID = DP.CARD_ID
            WHERE C.CATEGORY = 'RESPONSE'
              AND DP.LOBBY_ID = VAR_LOBBY_ID
            ORDER BY RAND()
            LIMIT 1;

            INSERT
            INTO HAND
                (
                    PLAYER_ID,
                    CARD_ID
                )
            VALUES
                (
                    VAR_PLAYER_ID,
                    VAR_CARD_ID
                );

            DELETE
            FROM DRAW_PILE
            WHERE LOBBY_ID = VAR_LOBBY_ID
              AND CARD_ID = VAR_CARD_ID;

            INSERT
            INTO LOG_DRAW
                (
                    LOBBY_ID,
                    USER_ID,
                    CARD_ID
                )
            VALUES
                (
                    VAR_LOBBY_ID,
                    VAR_PLAYER_USER_ID,
                    VAR_CARD_ID
                );

            SELECT
                COUNT(C.ID)
            INTO VAR_LOBBY_DRAW_PILE_SIZE
            FROM DRAW_PILE AS DP
                     INNER JOIN CARD AS C ON C.ID = DP.CARD_ID
            WHERE C.CATEGORY = 'RESPONSE'
              AND DP.LOBBY_ID = VAR_LOBBY_ID;

            SELECT
                COUNT(CARD_ID)
            INTO VAR_PLAYER_HAND_SIZE
            FROM HAND
            WHERE PLAYER_ID = VAR_PLAYER_ID;
        END WHILE;
END;