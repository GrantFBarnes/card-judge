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
            'Games Won' AS CATEGORY,
            1 AS GOAL
        UNION
        SELECT
            'WIN-GAME-10' AS CODE,
            'Games Won' AS CATEGORY,
            10 AS GOAL
        UNION
        SELECT
            'WIN-GAME-100' AS CODE,
            'Games Won' AS CATEGORY,
            100 AS GOAL
        UNION
        SELECT
            'WIN-ROUND-10' AS CODE,
            'Rounds Won' AS CATEGORY,
            10 AS GOAL
        UNION
        SELECT
            'WIN-ROUND-100' AS CODE,
            'Rounds Won' AS CATEGORY,
            100 AS GOAL
        UNION
        SELECT
            'WIN-ROUND-1000' AS CODE,
            'Rounds Won' AS CATEGORY,
            1000 AS GOAL
        UNION
        SELECT
            'GAMBLE-1' AS CODE,
            'Gambles Made' AS CATEGORY,
            1 AS GOAL
        UNION
        SELECT
            'GAMBLE-10' AS CODE,
            'Gambles Made' AS CATEGORY,
            10 AS GOAL
        UNION
        SELECT
            'GAMBLE-100' AS CODE,
            'Gambles Made' AS CATEGORY,
            100 AS GOAL
        UNION
        SELECT
            'BET-1' AS CODE,
            'Bets Placed' AS CATEGORY,
            1 AS GOAL
        UNION
        SELECT
            'BET-10' AS CODE,
            'Bets Placed' AS CATEGORY,
            10 AS GOAL
        UNION
        SELECT
            'BET-100' AS CODE,
            'Bets Placed' AS CATEGORY,
            100 AS GOAL
        UNION
        SELECT
            'PERK-1' AS CODE,
            'Perks Pruchased' AS CATEGORY,
            1 AS GOAL
        UNION
        SELECT
            'PERK-10' AS CODE,
            'Perks Pruchased' AS CATEGORY,
            10 AS GOAL
        UNION
        SELECT
            'PERK-100' AS CODE,
            'Perks Pruchased' AS CATEGORY,
            100 AS GOAL
        UNION
        SELECT
            'KICK-1' AS CODE,
            'Kicked From Lobby' AS CATEGORY,
            1 AS GOAL
        UNION
        SELECT
            'KICK-10' AS CODE,
            'Kicked From Lobby' AS CATEGORY,
            10 AS GOAL
        UNION
        SELECT
            'KICK-100' AS CODE,
            'Kicked From Lobby' AS CATEGORY,
            100 AS GOAL
    ) AS T;