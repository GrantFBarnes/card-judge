CREATE
OR REPLACE VIEW V_ACHIEVEMENT AS
SELECT
    CODE,
    NAME,
    ROWNUM () AS LIST_ORDER
FROM (
        SELECT
            'WIN-GAME-1' AS CODE,
            'Win a game' AS NAME
        UNION
        SELECT
            'WIN-GAME-10' AS CODE,
            'Win 10 games' AS NAME
        UNION
        SELECT
            'WIN-GAME-100' AS CODE,
            'Win 100 games' AS NAME
        UNION
        SELECT
            'WIN-ROUND-1' AS CODE,
            'Win a round' AS NAME
        UNION
        SELECT
            'WIN-ROUND-10' AS CODE,
            'Win 10 rounds' AS NAME
        UNION
        SELECT
            'WIN-ROUND-100' AS CODE,
            'Win 100 rounds' AS NAME
        UNION
        SELECT
            'WIN-ROUND-1000' AS CODE,
            'Win 1000 rounds' AS NAME
        UNION
        SELECT
            'GAMBLE-1' AS CODE,
            'Gamble 1 credit' AS NAME
        UNION
        SELECT
            'GAMBLE-10' AS CODE,
            'Gamble 10 credits' AS NAME
        UNION
        SELECT
            'GAMBLE-20' AS CODE,
            'Gamble 20 credits' AS NAME
        UNION
        SELECT
            'GAMBLE-WIN-2' AS CODE,
            'Win 2 credits gambling' AS NAME
        UNION
        SELECT
            'GAMBLE-WIN-20' AS CODE,
            'Win 20 credits gambling' AS NAME
        UNION
        SELECT
            'GAMBLE-WIN-40' AS CODE,
            'Win 40 credits gambling' AS NAME
        UNION
        SELECT
            'BET-1' AS CODE,
            'Bet 1 credit' AS NAME
        UNION
        SELECT
            'BET-10' AS CODE,
            'Bet 10 credits' AS NAME
        UNION
        SELECT
            'BET-20' AS CODE,
            'Bet 20 credits' AS NAME
        UNION
        SELECT
            'BET-WIN-2' AS CODE,
            'Win 2 credits betting' AS NAME
        UNION
        SELECT
            'BET-WIN-20' AS CODE,
            'Win 20 credits betting' AS NAME
        UNION
        SELECT
            'BET-WIN-40' AS CODE,
            'Win 40 credits betting' AS NAME
        UNION
        SELECT
            'KICK-1' AS CODE,
            'Get kicked' AS NAME
        UNION
        SELECT
            'KICK-10' AS CODE,
            'Get kicked 10 times' AS NAME
        UNION
        SELECT
            'KICK-20' AS CODE,
            'Get kicked 20 times' AS NAME
    ) AS T;