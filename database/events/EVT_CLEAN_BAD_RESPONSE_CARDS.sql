CREATE EVENT IF NOT EXISTS EVT_CLEAN_BAD_RESPONSE_CARDS ON SCHEDULE EVERY 1 DAY
DO
    BEGIN
        DECLARE VAR_DISCARD_COUNT INT;
        SET VAR_DISCARD_COUNT = 10;

        DELETE
        FROM CARD
        WHERE ID IN (
                SELECT
                    CARD_ID
                FROM (
                        SELECT
                            LD.CARD_ID,
                            COUNT(*) AS DISCARD_COUNT
                        FROM LOG_DISCARD AS LD
                            INNER JOIN CARD AS C ON C.ID = LD.CARD_ID
                            LEFT JOIN (
                                    SELECT
                                        PLAYER_CARD_ID AS CARD_ID,
                                        CREATED_ON_DATE AS LAST_PLAYED_DATE
                                    FROM (
                                            SELECT
                                                PLAYER_CARD_ID,
                                                CREATED_ON_DATE,
                                                RANK() OVER (
                                                    PARTITION BY PLAYER_CARD_ID
                                                    ORDER BY CREATED_ON_DATE DESC
                                                ) AS PLAY_ORDER
                                            FROM LOG_RESPONSE_CARD
                                        ) AS CARDSPLAYED
                                    WHERE PLAY_ORDER = 1
                                ) AS LASTPLAYED ON LASTPLAYED.CARD_ID = LD.CARD_ID
                        WHERE LASTPLAYED.LAST_PLAYED_DATE IS NULL
                            -- NEVER PLAYED
                            OR LASTPLAYED.LAST_PLAYED_DATE < LD.CREATED_ON_DATE
                        -- DISCARDS SINCE LAST PLAYED
                        GROUP BY LD.CARD_ID
                    ) AS BADCARDS
                WHERE DISCARD_COUNT > VAR_DISCARD_COUNT
            );
    END;