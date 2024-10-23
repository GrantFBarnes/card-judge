CREATE TABLE IF NOT EXISTS DRAW_PILE
(
    ID UUID NOT NULL DEFAULT UUID(),
    CREATED_ON_DATE DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(),

    LOBBY_ID UUID NOT NULL,
    CARD_ID UUID NOT NULL,

    PRIMARY KEY (ID),
    FOREIGN KEY (LOBBY_ID) REFERENCES LOBBY (ID) ON DELETE CASCADE,
    FOREIGN KEY (CARD_ID) REFERENCES CARD (ID) ON DELETE CASCADE,
    CONSTRAINT LOBBY_CARD_UNIQUE UNIQUE (LOBBY_ID, CARD_ID)
);