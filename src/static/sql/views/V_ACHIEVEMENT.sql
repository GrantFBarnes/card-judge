CREATE
OR REPLACE VIEW V_ACHIEVEMENT AS
SELECT
    CODE,
    CATEGORY,
    GOAL,
    ROWNUM () AS LIST_ORDER
FROM (
        SELECT
            'WIN-GAME-1' AS CODE,
            'Game Win' AS CATEGORY,
            1 AS GOAL
        UNION
        SELECT
            'WIN-GAME-10' AS CODE,
            'Game Win' AS CATEGORY,
            10 AS GOAL
        UNION
        SELECT
            'WIN-GAME-100' AS CODE,
            'Game Win' AS CATEGORY,
            100 AS GOAL
        UNION
        SELECT
            'WIN-ROUND-10' AS CODE,
            'Round Win' AS CATEGORY,
            10 AS GOAL
        UNION
        SELECT
            'WIN-ROUND-100' AS CODE,
            'Round Win' AS CATEGORY,
            100 AS GOAL
        UNION
        SELECT
            'WIN-ROUND-1000' AS CODE,
            'Round Win' AS CATEGORY,
            1000 AS GOAL
        UNION
        SELECT
            'GAMBLE-1' AS CODE,
            'Gamble' AS CATEGORY,
            1 AS GOAL
        UNION
        SELECT
            'GAMBLE-10' AS CODE,
            'Gamble' AS CATEGORY,
            10 AS GOAL
        UNION
        SELECT
            'GAMBLE-100' AS CODE,
            'Gamble' AS CATEGORY,
            100 AS GOAL
        UNION
        SELECT
            'BET-1' AS CODE,
            'Bet' AS CATEGORY,
            1 AS GOAL
        UNION
        SELECT
            'BET-10' AS CODE,
            'Bet' AS CATEGORY,
            10 AS GOAL
        UNION
        SELECT
            'BET-100' AS CODE,
            'Bet' AS CATEGORY,
            100 AS GOAL
        UNION
        SELECT
            'KICK-1' AS CODE,
            'Kicked' AS CATEGORY,
            1 AS GOAL
        UNION
        SELECT
            'KICK-10' AS CODE,
            'Kicked' AS CATEGORY,
            10 AS GOAL
        UNION
        SELECT
            'KICK-100' AS CODE,
            'Kicked' AS CATEGORY,
            100 AS GOAL
    ) AS T;