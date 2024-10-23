CREATE TABLE IF NOT EXISTS USER_ACCESS_LOBBY
(
    ID UUID NOT NULL DEFAULT UUID(),
    CREATED_ON_DATE DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(),

    USER_ID UUID NOT NULL,
    LOBBY_ID UUID NOT NULL,

    PRIMARY KEY (ID),
    FOREIGN KEY (USER_ID) REFERENCES USER (ID) ON DELETE CASCADE,
    FOREIGN KEY (LOBBY_ID) REFERENCES LOBBY (ID) ON DELETE CASCADE,
    CONSTRAINT USER_LOBBY_UNIQUE UNIQUE (USER_ID, LOBBY_ID)
);