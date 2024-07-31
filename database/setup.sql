DROP DATABASE IF EXISTS CARD_JUDGE;

CREATE DATABASE CARD_JUDGE
    CHARACTER SET = 'UTF8MB4'
    COLLATE = 'UTF8MB4_UNICODE_CI';

USE CARD_JUDGE;

CREATE TABLE PLAYER
(
    ID CHAR(36) NOT NULL DEFAULT UUID(),
    NAME VARCHAR(255) NOT NULL,
    DATE_ADDED DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(),
    DATE_MODIFIED DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(),
    PRIMARY KEY (ID),
    CONSTRAINT PLAYER_NAME_UNIQUE UNIQUE (NAME)
);

CREATE TABLE LOBBY
(
    ID CHAR(36) NOT NULL DEFAULT UUID(),
    NAME VARCHAR(255) NOT NULL,
    PASSWORD VARCHAR(255),
    DATE_ADDED DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(),
    DATE_MODIFIED DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(),
    PRIMARY KEY (ID),
    CONSTRAINT LOBBY_NAME_UNIQUE UNIQUE (NAME)
);

CREATE TABLE DECK
(
    ID CHAR(36) NOT NULL DEFAULT UUID(),
    NAME VARCHAR(255) NOT NULL,
    PASSWORD VARCHAR(255),
    DATE_ADDED DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(),
    DATE_MODIFIED DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(),
    PRIMARY KEY (ID),
    CONSTRAINT DECK_NAME_UNIQUE UNIQUE (NAME)
);

CREATE TABLE CARD
(
    ID CHAR(36) NOT NULL DEFAULT UUID(),
    DECK_ID CHAR(36) NOT NULL,
    TYPE ENUM ('Judge', 'Player') NOT NULL,
    TEXT VARCHAR(255) NOT NULL,
    DATE_ADDED DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(),
    DATE_MODIFIED DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(),
    PRIMARY KEY (ID),
    FOREIGN KEY (DECK_ID) REFERENCES DECK (ID) ON DELETE CASCADE,
    CONSTRAINT CARD_DECK_TEXT_UNIQUE UNIQUE (DECK_ID, TEXT)
);

CREATE TABLE LOBBY_DECK
(
    ID CHAR(36) NOT NULL DEFAULT UUID(),
    LOBBY_ID CHAR(36) NOT NULL,
    DECK_ID CHAR(36) NOT NULL,
    DATE_ADDED DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(),
    DATE_MODIFIED DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(),
    PRIMARY KEY (ID),
    FOREIGN KEY (LOBBY_ID) REFERENCES LOBBY (ID) ON DELETE CASCADE,
    FOREIGN KEY (DECK_ID) REFERENCES DECK (ID) ON DELETE CASCADE,
    CONSTRAINT LOBBY_DECK_UNIQUE UNIQUE (LOBBY_ID, DECK_ID)
);

CREATE TABLE LOBBY_PLAYER
(
    ID CHAR(36) NOT NULL DEFAULT UUID(),
    LOBBY_ID CHAR(36) NOT NULL,
    PLAYER_ID CHAR(36) NOT NULL,
    DATE_ADDED DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(),
    DATE_MODIFIED DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(),
    PRIMARY KEY (ID),
    FOREIGN KEY (LOBBY_ID) REFERENCES LOBBY (ID) ON DELETE CASCADE,
    FOREIGN KEY (PLAYER_ID) REFERENCES PLAYER (ID) ON DELETE CASCADE,
    CONSTRAINT LOBBY_PLAYER_UNIQUE UNIQUE (LOBBY_ID, PLAYER_ID)
);
