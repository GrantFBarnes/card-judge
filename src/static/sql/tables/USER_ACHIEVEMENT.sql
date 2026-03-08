CREATE TABLE IF NOT EXISTS USER_ACHIEVEMENT(
    USER_ID UUID NOT NULL,
    ACHIEVEMENT_CODE ENUM(
        'WIN-GAME-1',
        'WIN-GAME-10',
        'WIN-GAME-100',
        'WIN-ROUND-1',
        'WIN-ROUND-10',
        'WIN-ROUND-100',
        'WIN-ROUND-1000',
        'GAMBLE-1',
        'GAMBLE-10',
        'GAMBLE-20',
        'GAMBLE-WIN-2',
        'GAMBLE-WIN-20',
        'GAMBLE-WIN-40',
        'BET-1',
        'BET-10',
        'BET-20',
        'BET-WIN-2',
        'BET-WIN-20',
        'BET-WIN-40',
        'KICK-1',
        'KICK-10',
        'KICK-20'
    ) NOT NULL,
    PRIMARY KEY(USER_ID, ACHIEVEMENT_CODE),
    FOREIGN KEY(USER_ID) REFERENCES USER (ID) ON DELETE CASCADE
);