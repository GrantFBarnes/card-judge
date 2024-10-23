CREATE TABLE IF NOT EXISTS BOARD
(
    ID UUID NOT NULL DEFAULT UUID(),
    CREATED_ON_DATE DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(),

    LOBBY_ID UUID NOT NULL,
    PLAYER_ID UUID NOT NULL,
    CARD_ID UUID NOT NULL,
    SPECIAL_CATEGORY ENUM('STEAL','SURPRISE','WILD') NULL DEFAULT NULL,

    PRIMARY KEY (ID),
    FOREIGN KEY (LOBBY_ID) REFERENCES LOBBY (ID) ON DELETE CASCADE,
    FOREIGN KEY (PLAYER_ID) REFERENCES PLAYER (ID) ON DELETE CASCADE,
    FOREIGN KEY (CARD_ID) REFERENCES CARD (ID) ON DELETE CASCADE
);