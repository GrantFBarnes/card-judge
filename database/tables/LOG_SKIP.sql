CREATE TABLE IF NOT EXISTS LOG_SKIP
(
    ID UUID NOT NULL DEFAULT UUID(),
    CREATED_ON_DATE DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(),

    PLAYER_USER_ID UUID NULL,
    CARD_ID UUID NULL,

    PRIMARY KEY (ID),
    FOREIGN KEY (PLAYER_USER_ID) REFERENCES USER (ID) ON DELETE SET NULL,
    FOREIGN KEY (CARD_ID) REFERENCES CARD (ID) ON DELETE SET NULL
);