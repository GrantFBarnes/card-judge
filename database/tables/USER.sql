CREATE TABLE IF NOT EXISTS USER (
    ID UUID NOT NULL DEFAULT UUID(),
    CREATED_ON_DATE DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(),
    CHANGED_ON_DATE DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(),
    NAME VARCHAR(255) NOT NULL,
    PASSWORD_HASH CHAR(60) NOT NULL,
    COLOR_THEME VARCHAR(255) NULL,
    IS_APPROVED BOOLEAN NOT NULL DEFAULT 0,
    IS_ADMIN BOOLEAN NOT NULL DEFAULT 0,
    PRIMARY KEY(ID),
    CONSTRAINT NAME_UNIQUE UNIQUE(NAME)
);