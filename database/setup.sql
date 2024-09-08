DROP DATABASE IF EXISTS CARD_JUDGE;

CREATE DATABASE CARD_JUDGE
    CHARACTER SET = 'UTF8MB4'
    COLLATE = 'UTF8MB4_UNICODE_CI';

USE CARD_JUDGE;

CREATE TABLE PLAYER
(
    ID CHAR(36) NOT NULL DEFAULT UUID(),
    NAME VARCHAR(255) NOT NULL,
    PASSWORD_HASH CHAR(60) NOT NULL,
    COLOR_THEME VARCHAR(255) NULL,
    IS_ADMIN BOOLEAN NOT NULL DEFAULT 0,
    CREATED_ON_DATE DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(),
    CHANGED_ON_DATE DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(),
    PRIMARY KEY (ID),
    CONSTRAINT NAME_UNIQUE UNIQUE (NAME)
);

INSERT INTO PLAYER (ID, NAME, PASSWORD_HASH, IS_ADMIN)
VALUES ('ffb89249-6d88-11ef-aad4-28800dbd8d8a', 'Grant', '$2a$14$t7gWxR3Ak8uBkyPnw4TZz.WcN3nVlbDMEQgqHOuxEfWN3yCL3dgY.', 1);

CREATE TABLE DECK
(
    ID CHAR(36) NOT NULL DEFAULT UUID(),
    NAME VARCHAR(255) NOT NULL,
    PASSWORD_HASH CHAR(60) NULL,
    CREATED_ON_DATE DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(),
    CHANGED_ON_DATE DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(),
    CREATED_BY_PLAYER_ID CHAR(36) NOT NULL,
    CHANGED_BY_PLAYER_ID CHAR(36) NOT NULL,
    PRIMARY KEY (ID),
    CONSTRAINT NAME_UNIQUE UNIQUE (NAME)
);

INSERT INTO DECK (ID, NAME, PASSWORD_HASH, CREATED_BY_PLAYER_ID, CHANGED_BY_PLAYER_ID)
VALUES ('f395b797-6d89-11ef-aad4-28800dbd8d8a', 'Deck One', NULL, 'ffb89249-6d88-11ef-aad4-28800dbd8d8a', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a'),
       ('f395b862-6d89-11ef-aad4-28800dbd8d8a', 'Deck Two', NULL, 'ffb89249-6d88-11ef-aad4-28800dbd8d8a', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a'),
       ('f395b8e4-6d89-11ef-aad4-28800dbd8d8a', 'Deck Three', '$2a$14$t7gWxR3Ak8uBkyPnw4TZz.WcN3nVlbDMEQgqHOuxEfWN3yCL3dgY.', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a');

CREATE TABLE CARD
(
    ID CHAR(36) NOT NULL DEFAULT UUID(),
    DECK_ID CHAR(36) NOT NULL,
    TYPE ENUM ('Judge', 'Player') NOT NULL,
    TEXT VARCHAR(255) NOT NULL,
    CREATED_ON_DATE DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(),
    CHANGED_ON_DATE DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(),
    CREATED_BY_PLAYER_ID CHAR(36) NOT NULL,
    CHANGED_BY_PLAYER_ID CHAR(36) NOT NULL,
    PRIMARY KEY (ID),
    FOREIGN KEY (DECK_ID) REFERENCES DECK (ID) ON DELETE CASCADE,
    CONSTRAINT DECK_TEXT_UNIQUE UNIQUE (DECK_ID, TEXT)
);

INSERT INTO CARD (DECK_ID, TYPE, TEXT, CREATED_BY_PLAYER_ID, CHANGED_BY_PLAYER_ID)
VALUES ('f395b797-6d89-11ef-aad4-28800dbd8d8a', 'Judge', 'Deck One - Judge Card 1', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a'),
       ('f395b797-6d89-11ef-aad4-28800dbd8d8a', 'Judge', 'Deck One - Judge Card 2', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a'),
       ('f395b797-6d89-11ef-aad4-28800dbd8d8a', 'Judge', 'Deck One - Judge Card 3', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a'),
       ('f395b797-6d89-11ef-aad4-28800dbd8d8a', 'Judge', 'Deck One - Judge Card 4', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a'),
       ('f395b797-6d89-11ef-aad4-28800dbd8d8a', 'Judge', 'Deck One - Judge Card 5', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a'),
       ('f395b797-6d89-11ef-aad4-28800dbd8d8a', 'Player', 'Deck One - Player Card 1', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a'),
       ('f395b797-6d89-11ef-aad4-28800dbd8d8a', 'Player', 'Deck One - Player Card 2', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a'),
       ('f395b797-6d89-11ef-aad4-28800dbd8d8a', 'Player', 'Deck One - Player Card 3', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a'),
       ('f395b797-6d89-11ef-aad4-28800dbd8d8a', 'Player', 'Deck One - Player Card 4', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a'),
       ('f395b797-6d89-11ef-aad4-28800dbd8d8a', 'Player', 'Deck One - Player Card 5', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a'),
       ('f395b797-6d89-11ef-aad4-28800dbd8d8a', 'Player', 'Deck One - Player Card 6', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a'),
       ('f395b797-6d89-11ef-aad4-28800dbd8d8a', 'Player', 'Deck One - Player Card 7', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a'),
       ('f395b797-6d89-11ef-aad4-28800dbd8d8a', 'Player', 'Deck One - Player Card 8', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a'),
       ('f395b797-6d89-11ef-aad4-28800dbd8d8a', 'Player', 'Deck One - Player Card 9', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a'),
       ('f395b797-6d89-11ef-aad4-28800dbd8d8a', 'Player', 'Deck One - Player Card 10', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a');

INSERT INTO CARD (DECK_ID, TYPE, TEXT, CREATED_BY_PLAYER_ID, CHANGED_BY_PLAYER_ID)
VALUES ('f395b862-6d89-11ef-aad4-28800dbd8d8a', 'Judge', 'Deck Two - Judge Card 1', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a'),
       ('f395b862-6d89-11ef-aad4-28800dbd8d8a', 'Judge', 'Deck Two - Judge Card 2', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a'),
       ('f395b862-6d89-11ef-aad4-28800dbd8d8a', 'Judge', 'Deck Two - Judge Card 3', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a'),
       ('f395b862-6d89-11ef-aad4-28800dbd8d8a', 'Judge', 'Deck Two - Judge Card 4', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a'),
       ('f395b862-6d89-11ef-aad4-28800dbd8d8a', 'Judge', 'Deck Two - Judge Card 5', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a'),
       ('f395b862-6d89-11ef-aad4-28800dbd8d8a', 'Player', 'Deck Two - Player Card 1', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a'),
       ('f395b862-6d89-11ef-aad4-28800dbd8d8a', 'Player', 'Deck Two - Player Card 2', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a'),
       ('f395b862-6d89-11ef-aad4-28800dbd8d8a', 'Player', 'Deck Two - Player Card 3', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a'),
       ('f395b862-6d89-11ef-aad4-28800dbd8d8a', 'Player', 'Deck Two - Player Card 4', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a'),
       ('f395b862-6d89-11ef-aad4-28800dbd8d8a', 'Player', 'Deck Two - Player Card 5', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a'),
       ('f395b862-6d89-11ef-aad4-28800dbd8d8a', 'Player', 'Deck Two - Player Card 6', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a'),
       ('f395b862-6d89-11ef-aad4-28800dbd8d8a', 'Player', 'Deck Two - Player Card 7', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a'),
       ('f395b862-6d89-11ef-aad4-28800dbd8d8a', 'Player', 'Deck Two - Player Card 8', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a'),
       ('f395b862-6d89-11ef-aad4-28800dbd8d8a', 'Player', 'Deck Two - Player Card 9', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a'),
       ('f395b862-6d89-11ef-aad4-28800dbd8d8a', 'Player', 'Deck Two - Player Card 10', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a');

INSERT INTO CARD (DECK_ID, TYPE, TEXT, CREATED_BY_PLAYER_ID, CHANGED_BY_PLAYER_ID)
VALUES ('f395b8e4-6d89-11ef-aad4-28800dbd8d8a', 'Judge', 'Deck Three - Judge Card 1', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a'),
       ('f395b8e4-6d89-11ef-aad4-28800dbd8d8a', 'Judge', 'Deck Three - Judge Card 2', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a'),
       ('f395b8e4-6d89-11ef-aad4-28800dbd8d8a', 'Judge', 'Deck Three - Judge Card 3', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a'),
       ('f395b8e4-6d89-11ef-aad4-28800dbd8d8a', 'Judge', 'Deck Three - Judge Card 4', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a'),
       ('f395b8e4-6d89-11ef-aad4-28800dbd8d8a', 'Judge', 'Deck Three - Judge Card 5', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a'),
       ('f395b8e4-6d89-11ef-aad4-28800dbd8d8a', 'Player', 'Deck Three - Player Card 1', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a'),
       ('f395b8e4-6d89-11ef-aad4-28800dbd8d8a', 'Player', 'Deck Three - Player Card 2', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a'),
       ('f395b8e4-6d89-11ef-aad4-28800dbd8d8a', 'Player', 'Deck Three - Player Card 3', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a'),
       ('f395b8e4-6d89-11ef-aad4-28800dbd8d8a', 'Player', 'Deck Three - Player Card 4', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a'),
       ('f395b8e4-6d89-11ef-aad4-28800dbd8d8a', 'Player', 'Deck Three - Player Card 5', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a'),
       ('f395b8e4-6d89-11ef-aad4-28800dbd8d8a', 'Player', 'Deck Three - Player Card 6', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a'),
       ('f395b8e4-6d89-11ef-aad4-28800dbd8d8a', 'Player', 'Deck Three - Player Card 7', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a'),
       ('f395b8e4-6d89-11ef-aad4-28800dbd8d8a', 'Player', 'Deck Three - Player Card 8', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a'),
       ('f395b8e4-6d89-11ef-aad4-28800dbd8d8a', 'Player', 'Deck Three - Player Card 9', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a'),
       ('f395b8e4-6d89-11ef-aad4-28800dbd8d8a', 'Player', 'Deck Three - Player Card 10', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a', 'ffb89249-6d88-11ef-aad4-28800dbd8d8a');

CREATE TABLE LOBBY
(
    ID CHAR(36) NOT NULL DEFAULT UUID(),
    NAME VARCHAR(255) NOT NULL,
    PASSWORD_HASH CHAR(60) NULL,
    JUDGE_PLAYER_ID CHAR(36) NULL,
    JUDGE_CARD_ID CHAR(36) NULL,
    CREATED_ON_DATE DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(),
    CHANGED_ON_DATE DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(),
    CREATED_BY_PLAYER_ID CHAR(36) NOT NULL,
    CHANGED_BY_PLAYER_ID CHAR(36) NOT NULL,
    PRIMARY KEY (ID),
    FOREIGN KEY (JUDGE_PLAYER_ID) REFERENCES PLAYER (ID) ON DELETE SET NULL,
    FOREIGN KEY (JUDGE_CARD_ID) REFERENCES CARD (ID) ON DELETE SET NULL,
    CONSTRAINT NAME_UNIQUE UNIQUE (NAME)
);

CREATE TABLE PLAYER_ACCESS_DECK
(
    ID CHAR(36) NOT NULL DEFAULT UUID(),
    PLAYER_ID CHAR(36) NOT NULL,
    DECK_ID CHAR(36) NOT NULL,
    PRIMARY KEY (ID),
    FOREIGN KEY (PLAYER_ID) REFERENCES PLAYER (ID) ON DELETE CASCADE,
    FOREIGN KEY (DECK_ID) REFERENCES DECK (ID) ON DELETE CASCADE,
    CONSTRAINT PLAYER_DECK_UNIQUE UNIQUE (PLAYER_ID, DECK_ID)
);

CREATE TABLE PLAYER_ACCESS_LOBBY
(
    ID CHAR(36) NOT NULL DEFAULT UUID(),
    PLAYER_ID CHAR(36) NOT NULL,
    LOBBY_ID CHAR(36) NOT NULL,
    PRIMARY KEY (ID),
    FOREIGN KEY (PLAYER_ID) REFERENCES PLAYER (ID) ON DELETE CASCADE,
    FOREIGN KEY (LOBBY_ID) REFERENCES LOBBY (ID) ON DELETE CASCADE,
    CONSTRAINT PLAYER_LOBBY_UNIQUE UNIQUE (PLAYER_ID, LOBBY_ID)
);

CREATE TABLE LOBBY_CARD
(
    ID CHAR(36) NOT NULL DEFAULT UUID(),
    LOBBY_ID CHAR(36) NOT NULL,
    CARD_ID CHAR(36) NOT NULL,
    PRIMARY KEY (ID),
    FOREIGN KEY (LOBBY_ID) REFERENCES LOBBY (ID) ON DELETE CASCADE,
    FOREIGN KEY (CARD_ID) REFERENCES CARD (ID) ON DELETE CASCADE,
    CONSTRAINT LOBBY_CARD_UNIQUE UNIQUE (LOBBY_ID, CARD_ID)
);

CREATE TABLE LOBBY_PLAYER
(
    ID CHAR(36) NOT NULL DEFAULT UUID(),
    LOBBY_ID CHAR(36) NOT NULL,
    PLAYER_ID CHAR(36) NOT NULL,
    PRIMARY KEY (ID),
    FOREIGN KEY (LOBBY_ID) REFERENCES LOBBY (ID) ON DELETE CASCADE,
    FOREIGN KEY (PLAYER_ID) REFERENCES PLAYER (ID) ON DELETE CASCADE,
    CONSTRAINT LOBBY_PLAYER_UNIQUE UNIQUE (LOBBY_ID, PLAYER_ID)
);

CREATE TABLE LOBBY_PLAYER_CARD
(
    ID CHAR(36) NOT NULL DEFAULT UUID(),
    LOBBY_PLAYER_ID CHAR(36) NOT NULL,
    LOBBY_ID CHAR(36) NOT NULL,
    PLAYER_ID CHAR(36) NOT NULL,
    CARD_ID CHAR(36) NOT NULL,
    PRIMARY KEY (ID),
    FOREIGN KEY (LOBBY_PLAYER_ID) REFERENCES LOBBY_PLAYER (ID) ON DELETE CASCADE,
    FOREIGN KEY (LOBBY_ID) REFERENCES LOBBY (ID) ON DELETE CASCADE,
    FOREIGN KEY (PLAYER_ID) REFERENCES PLAYER (ID) ON DELETE CASCADE,
    FOREIGN KEY (CARD_ID) REFERENCES CARD (ID) ON DELETE CASCADE,
    CONSTRAINT LOBBY_PLAYER_CARD_UNIQUE UNIQUE (LOBBY_PLAYER_ID, CARD_ID)
);

CREATE TABLE LOBBY_PLAYER_PLAY
(
    ID CHAR(36) NOT NULL DEFAULT UUID(),
    LOBBY_PLAYER_ID CHAR(36) NOT NULL,
    PLAY_CARD_ID CHAR(36) NOT NULL,
    PRIMARY KEY (ID),
    FOREIGN KEY (LOBBY_PLAYER_ID) REFERENCES LOBBY_PLAYER (ID) ON DELETE CASCADE,
    FOREIGN KEY (PLAY_CARD_ID) REFERENCES CARD (ID) ON DELETE CASCADE,
    CONSTRAINT LOBBY_PLAYER_UNIQUE UNIQUE (LOBBY_PLAYER_ID)
);

CREATE TABLE LOBBY_RESULT
(
    ID CHAR(36) NOT NULL DEFAULT UUID(),
    LOBBY_PLAYER_ID CHAR(36) NOT NULL,
    JUDGE_PLAYER_ID CHAR(36) NOT NULL,
    PRIMARY KEY (ID),
    FOREIGN KEY (LOBBY_PLAYER_ID) REFERENCES LOBBY_PLAYER (ID) ON DELETE CASCADE,
    FOREIGN KEY (JUDGE_PLAYER_ID) REFERENCES PLAYER (ID) ON DELETE CASCADE
);
