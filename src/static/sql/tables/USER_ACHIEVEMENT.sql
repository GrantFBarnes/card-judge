CREATE TABLE IF NOT EXISTS USER_ACHIEVEMENT(
    USER_ID UUID NOT NULL,
    ACHIEVEMENT_CODE ENUM(
        'WIN-GAME-1',
        'WIN-GAME-10',
        'WIN-GAME-100',
        'WIN-ROUND-10',
        'WIN-ROUND-100',
        'WIN-ROUND-1000',
        'GAMBLE-1',
        'GAMBLE-10',
        'GAMBLE-100',
        'BET-1',
        'BET-10',
        'BET-100',
        'PERK-1',
        'PERK-10',
        'PERK-100',
        'KICK-1',
        'KICK-10',
        'KICK-100',
        'FLIP-TABLE-1',
        'FLIP-TABLE-10',
        'FLIP-TABLE-100'
    ) NOT NULL,
    PRIMARY KEY(USER_ID, ACHIEVEMENT_CODE),
    FOREIGN KEY(USER_ID) REFERENCES USER (ID) ON DELETE CASCADE
);