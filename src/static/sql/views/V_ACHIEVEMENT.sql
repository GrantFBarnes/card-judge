CREATE
OR REPLACE VIEW V_ACHIEVEMENT AS
SELECT
    CODE,
    CATEGORY,
    THRESHOLD,
    ROWNUM () AS LIST_ORDER
FROM (
        SELECT
            'WIN-GAME-1' AS CODE,
            'Game Win' AS CATEGORY,
            1 AS THRESHOLD
        UNION
        SELECT
            'WIN-GAME-10' AS CODE,
            'Game Win' AS CATEGORY,
            10 AS THRESHOLD
        UNION
        SELECT
            'WIN-GAME-100' AS CODE,
            'Game Win' AS CATEGORY,
            100 AS THRESHOLD
        UNION
        SELECT
            'WIN-ROUND-10' AS CODE,
            'Round Win' AS CATEGORY,
            10 AS THRESHOLD
        UNION
        SELECT
            'WIN-ROUND-100' AS CODE,
            'Round Win' AS CATEGORY,
            100 AS THRESHOLD
        UNION
        SELECT
            'WIN-ROUND-1000' AS CODE,
            'Round Win' AS CATEGORY,
            1000 AS THRESHOLD
        UNION
        SELECT
            'GAMBLE-1' AS CODE,
            'Gamble' AS CATEGORY,
            1 AS THRESHOLD
        UNION
        SELECT
            'GAMBLE-10' AS CODE,
            'Gamble' AS CATEGORY,
            10 AS THRESHOLD
        UNION
        SELECT
            'GAMBLE-50' AS CODE,
            'Gamble' AS CATEGORY,
            50 AS THRESHOLD
        UNION
        SELECT
            'GAMBLE-WIN-2' AS CODE,
            'Gamble Win' AS CATEGORY,
            2 AS THRESHOLD
        UNION
        SELECT
            'GAMBLE-WIN-20' AS CODE,
            'Gamble Win' AS CATEGORY,
            20 AS THRESHOLD
        UNION
        SELECT
            'GAMBLE-WIN-100' AS CODE,
            'Gamble Win' AS CATEGORY,
            100 AS THRESHOLD
        UNION
        SELECT
            'BET-1' AS CODE,
            'Bet' AS CATEGORY,
            1 AS THRESHOLD
        UNION
        SELECT
            'BET-10' AS CODE,
            'Bet' AS CATEGORY,
            10 AS THRESHOLD
        UNION
        SELECT
            'BET-50' AS CODE,
            'Bet' AS CATEGORY,
            50 AS THRESHOLD
        UNION
        SELECT
            'BET-WIN-2' AS CODE,
            'Bet Win' AS CATEGORY,
            2 AS THRESHOLD
        UNION
        SELECT
            'BET-WIN-20' AS CODE,
            'Bet Win' AS CATEGORY,
            20 AS THRESHOLD
        UNION
        SELECT
            'BET-WIN-100' AS CODE,
            'Bet Win' AS CATEGORY,
            100 AS THRESHOLD
        UNION
        SELECT
            'KICK-1' AS CODE,
            'Kicked' AS CATEGORY,
            1 AS THRESHOLD
        UNION
        SELECT
            'KICK-10' AS CODE,
            'Kicked' AS CATEGORY,
            10 AS THRESHOLD
        UNION
        SELECT
            'KICK-50' AS CODE,
            'Kicked' AS CATEGORY,
            50 AS THRESHOLD
    ) AS T;