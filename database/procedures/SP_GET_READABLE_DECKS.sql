CREATE PROCEDURE IF NOT EXISTS SP_GET_READABLE_DECKS (
    IN VAR_USER_ID UUID
)
BEGIN
    IF EXISTS(SELECT ID
        FROM USER
        WHERE ID = VAR_USER_ID
            AND IS_ADMIN = 1) THEN
        SELECT
            D.ID,
            D.NAME
        FROM DECK AS D
        WHERE D.IS_LOBBY_WILD_DECK = FALSE
        ORDER BY D.NAME;
        ELSE
        SELECT
            DISTINCT
            D.ID,
            D.NAME
        FROM DECK AS D
            LEFT JOIN USER_ACCESS_DECK AS UAD ON UAD.DECK_ID = D.ID
        WHERE D.IS_LOBBY_WILD_DECK = FALSE
            AND (
                UAD.USER_ID = VAR_USER_ID
                OR D.IS_PUBLIC_READONLY
            )
        ORDER BY D.NAME;
    END IF;
END;