CREATE TABLE IF NOT EXISTS LOBBY
(
    ID UUID NOT NULL DEFAULT UUID(),
    CREATED_ON_DATE DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(),
    CHANGED_ON_DATE DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(),

    NAME VARCHAR(255) NOT NULL,
    PASSWORD_HASH CHAR(60) NULL,
    HAND_SIZE INT NOT NULL DEFAULT 8,
    CREDIT_LIMIT INT NOT NULL DEFAULT 3,

    PRIMARY KEY (ID),
    CONSTRAINT NAME_UNIQUE UNIQUE (NAME)
);