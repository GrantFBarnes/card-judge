DROP DATABASE IF EXISTS CARD_JUDGE;

CREATE DATABASE CARD_JUDGE
    CHARACTER SET = 'UTF8MB4'
    COLLATE = 'UTF8MB4_UNICODE_CI';

USE CARD_JUDGE;

CREATE TABLE USER
(
    ID CHAR(36) NOT NULL DEFAULT UUID(),
    CREATED_ON_DATE DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(),
    CHANGED_ON_DATE DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(),

    NAME VARCHAR(255) NOT NULL,
    PASSWORD_HASH CHAR(60) NOT NULL,
    COLOR_THEME VARCHAR(255) NULL,
    IS_ADMIN BOOLEAN NOT NULL DEFAULT 0,

    PRIMARY KEY (ID),
    CONSTRAINT NAME_UNIQUE UNIQUE (NAME)
);

CREATE TRIGGER TR_USER_BF_UP_SET_CHANGED_ON_DATE
BEFORE UPDATE ON USER
FOR EACH ROW
SET NEW.CHANGED_ON_DATE = CURRENT_TIMESTAMP();

CREATE TABLE DECK
(
    ID CHAR(36) NOT NULL DEFAULT UUID(),
    CREATED_ON_DATE DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(),
    CHANGED_ON_DATE DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(),

    NAME VARCHAR(255) NOT NULL,
    PASSWORD_HASH CHAR(60) NULL,

    PRIMARY KEY (ID),
    CONSTRAINT NAME_UNIQUE UNIQUE (NAME)
);

CREATE TRIGGER TR_DECK_BF_UP_SET_CHANGED_ON_DATE
BEFORE UPDATE ON DECK
FOR EACH ROW
SET NEW.CHANGED_ON_DATE = CURRENT_TIMESTAMP();

CREATE TABLE CARD_TYPE
(
    ID CHAR(36) NOT NULL DEFAULT UUID(),
    CREATED_ON_DATE DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(),
    CHANGED_ON_DATE DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(),

    NAME NVARCHAR(255) NOT NULL,

    PRIMARY KEY (ID),
    CONSTRAINT NAME_UNIQUE UNIQUE (NAME)
);

CREATE TRIGGER TR_CARD_TYPE_BF_UP_SET_CHANGED_ON_DATE
BEFORE UPDATE ON CARD_TYPE
FOR EACH ROW
SET NEW.CHANGED_ON_DATE = CURRENT_TIMESTAMP();

CREATE TABLE CARD
(
    ID CHAR(36) NOT NULL DEFAULT UUID(),
    CREATED_ON_DATE DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(),
    CHANGED_ON_DATE DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(),

    DECK_ID CHAR(36) NOT NULL,
    CARD_TYPE_ID CHAR(36) NOT NULL,
    TEXT VARCHAR(510) NOT NULL,
    BLANK_COUNT INT NOT NULL DEFAULT 0,
    SKIP_COUNT INT NOT NULL DEFAULT 0,

    PRIMARY KEY (ID),
    FOREIGN KEY (DECK_ID) REFERENCES DECK (ID) ON DELETE CASCADE,
    FOREIGN KEY (CARD_TYPE_ID) REFERENCES CARD_TYPE (ID) ON DELETE RESTRICT,
    CONSTRAINT DECK_TEXT_UNIQUE UNIQUE (DECK_ID, TEXT, CARD_TYPE_ID)
);

CREATE TRIGGER TR_CARD_BF_UP_SET_CHANGED_ON_DATE
BEFORE UPDATE ON CARD
FOR EACH ROW
SET NEW.CHANGED_ON_DATE = CURRENT_TIMESTAMP();

CREATE TABLE LOBBY
(
    ID CHAR(36) NOT NULL DEFAULT UUID(),
    CREATED_ON_DATE DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(),
    CHANGED_ON_DATE DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(),

    NAME VARCHAR(255) NOT NULL,
    PASSWORD_HASH CHAR(60) NULL,
    HAND_SIZE INT NOT NULL DEFAULT 8,

    PRIMARY KEY (ID),
    CONSTRAINT NAME_UNIQUE UNIQUE (NAME)
);

CREATE TRIGGER TR_LOBBY_BF_UP_SET_CHANGED_ON_DATE
BEFORE UPDATE ON LOBBY
FOR EACH ROW
SET NEW.CHANGED_ON_DATE = CURRENT_TIMESTAMP();

CREATE TABLE USER_ACCESS_DECK
(
    ID CHAR(36) NOT NULL DEFAULT UUID(),
    CREATED_ON_DATE DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(),
    CHANGED_ON_DATE DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(),

    USER_ID CHAR(36) NOT NULL,
    DECK_ID CHAR(36) NOT NULL,

    PRIMARY KEY (ID),
    FOREIGN KEY (USER_ID) REFERENCES USER (ID) ON DELETE CASCADE,
    FOREIGN KEY (DECK_ID) REFERENCES DECK (ID) ON DELETE CASCADE,
    CONSTRAINT USER_DECK_UNIQUE UNIQUE (USER_ID, DECK_ID)
);

CREATE TRIGGER TR_USER_ACCESS_DECK_BF_UP_SET_CHANGED_ON_DATE
BEFORE UPDATE ON USER_ACCESS_DECK
FOR EACH ROW
SET NEW.CHANGED_ON_DATE = CURRENT_TIMESTAMP();

DELIMITER //
CREATE TRIGGER TR_DECK_AF_UP_REVOKE_ACCESS
AFTER UPDATE ON DECK
FOR EACH ROW
BEGIN
    IF OLD.PASSWORD_HASH <> NEW.PASSWORD_HASH THEN
        DELETE FROM USER_ACCESS_DECK
        WHERE DECK_ID = NEW.ID;
    END IF;
END;
//
DELIMITER ;

CREATE TABLE USER_ACCESS_LOBBY
(
    ID CHAR(36) NOT NULL DEFAULT UUID(),
    CREATED_ON_DATE DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(),
    CHANGED_ON_DATE DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(),

    USER_ID CHAR(36) NOT NULL,
    LOBBY_ID CHAR(36) NOT NULL,

    PRIMARY KEY (ID),
    FOREIGN KEY (USER_ID) REFERENCES USER (ID) ON DELETE CASCADE,
    FOREIGN KEY (LOBBY_ID) REFERENCES LOBBY (ID) ON DELETE CASCADE,
    CONSTRAINT USER_LOBBY_UNIQUE UNIQUE (USER_ID, LOBBY_ID)
);

CREATE TRIGGER TR_USER_ACCESS_LOBBY_BF_UP_SET_CHANGED_ON_DATE
BEFORE UPDATE ON USER_ACCESS_LOBBY
FOR EACH ROW
SET NEW.CHANGED_ON_DATE = CURRENT_TIMESTAMP();

CREATE TABLE DRAW_PILE
(
    ID CHAR(36) NOT NULL DEFAULT UUID(),
    CREATED_ON_DATE DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(),
    CHANGED_ON_DATE DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(),

    LOBBY_ID CHAR(36) NOT NULL,
    CARD_ID CHAR(36) NOT NULL,

    PRIMARY KEY (ID),
    FOREIGN KEY (LOBBY_ID) REFERENCES LOBBY (ID) ON DELETE CASCADE,
    FOREIGN KEY (CARD_ID) REFERENCES CARD (ID) ON DELETE CASCADE,
    CONSTRAINT LOBBY_CARD_UNIQUE UNIQUE (LOBBY_ID, CARD_ID)
);

CREATE TRIGGER TR_DRAW_PILE_BF_UP_SET_CHANGED_ON_DATE
BEFORE UPDATE ON DRAW_PILE
FOR EACH ROW
SET NEW.CHANGED_ON_DATE = CURRENT_TIMESTAMP();

CREATE TABLE PLAYER
(
    ID CHAR(36) NOT NULL DEFAULT UUID(),
    CREATED_ON_DATE DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(),
    CHANGED_ON_DATE DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(),

    LOBBY_ID CHAR(36) NOT NULL,
    USER_ID CHAR(36) NOT NULL,
    ACTIVE BOOLEAN NOT NULL DEFAULT TRUE,

    PRIMARY KEY (ID),
    FOREIGN KEY (LOBBY_ID) REFERENCES LOBBY (ID) ON DELETE CASCADE,
    FOREIGN KEY (USER_ID) REFERENCES USER (ID) ON DELETE CASCADE,
    CONSTRAINT LOBBY_USER_UNIQUE UNIQUE (LOBBY_ID, USER_ID)
);

CREATE TRIGGER TR_PLAYER_BF_UP_SET_CHANGED_ON_DATE
BEFORE UPDATE ON PLAYER
FOR EACH ROW
SET NEW.CHANGED_ON_DATE = CURRENT_TIMESTAMP();

DELIMITER //
CREATE PROCEDURE SP_ADD_PLAYER (
    IN VAR_PLAYER_ID UUID,
    IN VAR_LOBBY_ID UUID,
    IN VAR_USER_ID UUID
)
BEGIN
    INSERT IGNORE INTO PLAYER (ID, LOBBY_ID, USER_ID)
    VALUES (VAR_PLAYER_ID, VAR_LOBBY_ID, VAR_USER_ID);

    UPDATE PLAYER
    SET ACTIVE = 1
    WHERE ID = VAR_PLAYER_ID;
END;
//
DELIMITER ;

CREATE TABLE HAND
(
    ID CHAR(36) NOT NULL DEFAULT UUID(),
    CREATED_ON_DATE DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(),
    CHANGED_ON_DATE DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(),

    PLAYER_ID CHAR(36) NOT NULL,
    CARD_ID CHAR(36) NOT NULL,

    PRIMARY KEY (ID),
    FOREIGN KEY (PLAYER_ID) REFERENCES PLAYER (ID) ON DELETE CASCADE,
    FOREIGN KEY (CARD_ID) REFERENCES CARD (ID) ON DELETE CASCADE,
    CONSTRAINT PLAYER_CARD_UNIQUE UNIQUE (PLAYER_ID, CARD_ID)
);

CREATE TRIGGER TR_HAND_BF_UP_SET_CHANGED_ON_DATE
BEFORE UPDATE ON HAND
FOR EACH ROW
SET NEW.CHANGED_ON_DATE = CURRENT_TIMESTAMP();

DELIMITER //
CREATE PROCEDURE SP_DRAW_HAND (
    IN HAND_PLAYER_ID UUID
)
BEGIN
    DECLARE LOBBY_HAND_SIZE INT;
    DECLARE PLAYER_HAND_SIZE INT;
    DECLARE CARDS_TO_DRAW INT;

    SELECT L.HAND_SIZE
    INTO LOBBY_HAND_SIZE
    FROM LOBBY AS L
        INNER JOIN PLAYER AS P ON P.LOBBY_ID = L.ID
    WHERE P.ID = HAND_PLAYER_ID;

    SELECT COUNT(CARD_ID)
    INTO PLAYER_HAND_SIZE
    FROM HAND
    WHERE PLAYER_ID = HAND_PLAYER_ID;

    SET CARDS_TO_DRAW = LOBBY_HAND_SIZE - PLAYER_HAND_SIZE;

    INSERT INTO HAND (
        PLAYER_ID,
        CARD_ID
    )
    SELECT
        P.ID AS PLAYER_ID,
        C.ID AS CARD_ID
    FROM DRAW_PILE AS DP
        INNER JOIN PLAYER AS P ON P.LOBBY_ID = DP.LOBBY_ID
        INNER JOIN CARD AS C ON C.ID = DP.CARD_ID
        INNER JOIN CARD_TYPE AS CT ON CT.ID = C.CARD_TYPE_ID
    WHERE CT.NAME = 'PLAYER'
        AND P.ID = HAND_PLAYER_ID
    ORDER BY RAND()
    LIMIT CARDS_TO_DRAW;

    DELETE DP
    FROM DRAW_PILE AS DP
        INNER JOIN PLAYER AS P ON P.LOBBY_ID = DP.LOBBY_ID
        INNER JOIN HAND AS H ON H.PLAYER_ID = P.ID AND H.CARD_ID = DP.CARD_ID
    WHERE P.ID = HAND_PLAYER_ID;

    INSERT INTO LOG_DRAW (
        PLAYER_USER_ID,
        CARD_ID
    )
    SELECT
        P.USER_ID AS PLAYER_USER_ID,
        H.CARD_ID
    FROM HAND AS H
        INNER JOIN PLAYER AS P ON P.ID = H.PLAYER_ID
    WHERE P.ID = HAND_PLAYER_ID
    ORDER BY H.CREATED_ON_DATE DESC
    LIMIT CARDS_TO_DRAW;
END;
//
DELIMITER ;

CREATE TRIGGER TR_PLAYER_AF_IN_DRAW_HAND
AFTER INSERT ON PLAYER
FOR EACH ROW
CALL SP_DRAW_HAND (NEW.ID);

CREATE TABLE JUDGE
(
    ID CHAR(36) NOT NULL DEFAULT UUID(),
    CREATED_ON_DATE DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(),
    CHANGED_ON_DATE DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(),

    LOBBY_ID CHAR(36) NOT NULL,
    PLAYER_ID CHAR(36) NOT NULL,
    CARD_ID CHAR(36) NOT NULL,

    PRIMARY KEY (ID),
    FOREIGN KEY (LOBBY_ID) REFERENCES LOBBY (ID) ON DELETE CASCADE,
    FOREIGN KEY (PLAYER_ID) REFERENCES PLAYER (ID) ON DELETE CASCADE,
    FOREIGN KEY (CARD_ID) REFERENCES CARD (ID) ON DELETE CASCADE,
    CONSTRAINT LOBBY_UNIQUE UNIQUE (LOBBY_ID),
    CONSTRAINT PLAYER_UNIQUE UNIQUE (PLAYER_ID)
);

CREATE TRIGGER TR_JUDGE_BF_UP_SET_CHANGED_ON_DATE
BEFORE UPDATE ON JUDGE
FOR EACH ROW
SET NEW.CHANGED_ON_DATE = CURRENT_TIMESTAMP();

DELIMITER //
CREATE PROCEDURE SP_SET_JUDGE (
    IN JUDGE_LOBBY_ID UUID
)
BEGIN
    INSERT IGNORE INTO JUDGE (
        LOBBY_ID,
        PLAYER_ID,
        CARD_ID
    )
    SELECT
        JUDGE_LOBBY_ID AS LOBBY_ID,
        P.ID AS PLAYER_ID,
        C.ID AS CARD_ID
    FROM DRAW_PILE AS DP
        INNER JOIN PLAYER AS P ON P.LOBBY_ID = DP.LOBBY_ID
        INNER JOIN CARD AS C ON C.ID = DP.CARD_ID
        INNER JOIN CARD_TYPE AS CT ON CT.ID = C.CARD_TYPE_ID
    WHERE CT.NAME = 'JUDGE'
        AND DP.LOBBY_ID = JUDGE_LOBBY_ID
    ORDER BY RAND()
    LIMIT 1;

    DELETE DP
    FROM DRAW_PILE AS DP
        INNER JOIN JUDGE AS J ON J.LOBBY_ID = DP.LOBBY_ID AND J.CARD_ID = DP.CARD_ID;
END;
//
DELIMITER ;

CREATE TRIGGER TR_PLAYER_AF_IN_SET_JUDGE
AFTER INSERT ON PLAYER
FOR EACH ROW
CALL SP_SET_JUDGE (NEW.LOBBY_ID);

CREATE TRIGGER TR_PLAYER_AF_DL_SET_JUDGE
AFTER DELETE ON PLAYER
FOR EACH ROW
CALL SP_SET_JUDGE (OLD.LOBBY_ID);

DELIMITER //
CREATE PROCEDURE SP_SKIP_JUDGE_CARD (
    IN JUDGE_LOBBY_ID UUID
)
BEGIN
    UPDATE CARD
    SET SKIP_COUNT = SKIP_COUNT + 1
    WHERE ID = (
        SELECT CARD_ID
        FROM JUDGE
        WHERE LOBBY_ID = JUDGE_LOBBY_ID
    );

    UPDATE JUDGE
    SET
        CARD_ID = (
            SELECT C.ID
            FROM DRAW_PILE AS DP
                INNER JOIN CARD AS C ON C.ID = DP.CARD_ID
                INNER JOIN CARD_TYPE AS CT ON CT.ID = C.CARD_TYPE_ID
            WHERE CT.NAME = 'JUDGE'
                AND DP.LOBBY_ID = JUDGE_LOBBY_ID
            ORDER BY RAND()
            LIMIT 1
        )
    WHERE LOBBY_ID = JUDGE_LOBBY_ID;

    DELETE DP
    FROM DRAW_PILE AS DP
        INNER JOIN JUDGE AS J ON J.LOBBY_ID = DP.LOBBY_ID AND J.CARD_ID = DP.CARD_ID;
END;
//
DELIMITER ;

CREATE TABLE BOARD
(
    ID CHAR(36) NOT NULL DEFAULT UUID(),
    CREATED_ON_DATE DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(),
    CHANGED_ON_DATE DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(),

    LOBBY_ID CHAR(36) NOT NULL,
    PLAYER_ID CHAR(36) NOT NULL,
    CARD_ID CHAR(36) NOT NULL,

    PRIMARY KEY (ID),
    FOREIGN KEY (LOBBY_ID) REFERENCES LOBBY (ID) ON DELETE CASCADE,
    FOREIGN KEY (PLAYER_ID) REFERENCES PLAYER (ID) ON DELETE CASCADE,
    FOREIGN KEY (CARD_ID) REFERENCES CARD (ID) ON DELETE CASCADE,
    CONSTRAINT PLAYER_UNIQUE UNIQUE (PLAYER_ID)
);

CREATE TRIGGER TR_BOARD_BF_UP_SET_CHANGED_ON_DATE
BEFORE UPDATE ON BOARD
FOR EACH ROW
SET NEW.CHANGED_ON_DATE = CURRENT_TIMESTAMP();

CREATE TABLE WIN
(
    ID CHAR(36) NOT NULL DEFAULT UUID(),
    CREATED_ON_DATE DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(),
    CHANGED_ON_DATE DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(),

    LOBBY_ID CHAR(36) NOT NULL,
    PLAYER_ID CHAR(36) NOT NULL,

    PRIMARY KEY (ID),
    FOREIGN KEY (LOBBY_ID) REFERENCES LOBBY (ID) ON DELETE CASCADE,
    FOREIGN KEY (PLAYER_ID) REFERENCES PLAYER (ID) ON DELETE CASCADE
);

CREATE TRIGGER TR_WIN_BF_UP_SET_CHANGED_ON_DATE
BEFORE UPDATE ON WIN
FOR EACH ROW
SET NEW.CHANGED_ON_DATE = CURRENT_TIMESTAMP();

DELIMITER //
CREATE PROCEDURE SP_REPLACE_JUDGE (
    IN JUDGE_LOBBY_ID UUID,
    IN NEW_PLAYER_ID UUID
)
BEGIN
    UPDATE JUDGE
    SET
        PLAYER_ID = NEW_PLAYER_ID,
        CARD_ID = (
            SELECT C.ID
            FROM DRAW_PILE AS DP
                INNER JOIN CARD AS C ON C.ID = DP.CARD_ID
                INNER JOIN CARD_TYPE AS CT ON CT.ID = C.CARD_TYPE_ID
            WHERE CT.NAME = 'JUDGE'
                AND DP.LOBBY_ID = JUDGE_LOBBY_ID
            ORDER BY RAND()
            LIMIT 1
        )
    WHERE LOBBY_ID = JUDGE_LOBBY_ID;

    DELETE DP
    FROM DRAW_PILE AS DP
        INNER JOIN JUDGE AS J ON J.LOBBY_ID = DP.LOBBY_ID AND J.CARD_ID = DP.CARD_ID;
END;
//
DELIMITER ;

CREATE TRIGGER TR_WIN_AF_IN_REPLACE_JUDGE
AFTER INSERT ON WIN
FOR EACH ROW
CALL SP_REPLACE_JUDGE (NEW.LOBBY_ID, NEW.PLAYER_ID);

DELIMITER //
CREATE PROCEDURE SP_PLAY_CARD (
    IN PLAYER_ID UUID,
    IN CARD_ID UUID
)
BEGIN
    DECLARE LOBBY_ID UUID;
    DECLARE PLAYER_USER_ID UUID;
    DECLARE JUDGE_USER_ID UUID;

    SELECT LOBBY_ID
    INTO LOBBY_ID
    FROM PLAYER
    WHERE ID = PLAYER_ID;

    SELECT USER_ID
    INTO PLAYER_USER_ID
    FROM PLAYER
    WHERE ID = PLAYER_ID;

    SELECT P.USER_ID
    INTO JUDGE_USER_ID
    FROM JUDGE AS J
        INNER JOIN PLAYER AS P ON P.ID = J.PLAYER_ID
    WHERE J.LOBBY_ID = LOBBY_ID;

    INSERT INTO BOARD (LOBBY_ID, PLAYER_ID, CARD_ID)
    VALUES (LOBBY_ID, PLAYER_ID, CARD_ID);

    INSERT INTO LOG_PLAY (PLAYER_USER_ID, JUDGE_USER_ID, CARD_ID)
    VALUES(PLAYER_USER_ID, JUDGE_USER_ID, CARD_ID);

    DELETE FROM HAND
    WHERE PLAYER_ID = PLAYER_ID
        AND CARD_ID = CARD_ID;
END;
//
DELIMITER ;

DELIMITER //
CREATE PROCEDURE SP_DISCARD_HAND (
    IN PLAYER_ID UUID
)
BEGIN
    DECLARE PLAYER_USER_ID UUID;

    SELECT USER_ID
    INTO PLAYER_USER_ID
    FROM PLAYER
    WHERE ID = PLAYER_ID;

    INSERT INTO LOG_DISCARD (PLAYER_USER_ID, CARD_ID)
    SELECT
        USER_ID AS PLAYER_USER_ID,
        CARD_ID
    FROM HAND
    WHERE PLAYER_ID = PLAYER_ID;

    DELETE FROM HAND
    WHERE PLAYER_ID = PLAYER_ID;
END;
//
DELIMITER ;

DELIMITER //
CREATE PROCEDURE SP_DISCARD_CARD (
    IN PLAYER_ID UUID,
    IN CARD_ID UUID
)
BEGIN
    DECLARE PLAYER_USER_ID UUID;

    SELECT USER_ID
    INTO PLAYER_USER_ID
    FROM PLAYER
    WHERE ID = PLAYER_ID;

    DELETE FROM HAND
    WHERE PLAYER_ID = PLAYER_ID
        AND CARD_ID = CARD_ID;

    INSERT INTO LOG_DISCARD (PLAYER_USER_ID, CARD_ID)
    VALUES(PLAYER_USER_ID, CARD_ID);
END;
//
DELIMITER ;

DELIMITER //
CREATE PROCEDURE SP_PICK_WINNER (
    IN WIN_LOBBY_ID UUID,
    IN WIN_CARD_ID UUID
)
BEGIN
    DECLARE WIN_PLAYER_ID UUID;
    DECLARE WIN_PLAYER_USER_ID UUID;
    DECLARE WIN_JUDGE_USER_ID UUID;

    SELECT PLAYER_ID
    INTO WIN_PLAYER_ID
    FROM BOARD
    WHERE LOBBY_ID = WIN_LOBBY_ID
        AND CARD_ID = WIN_CARD_ID;

    SELECT USER_ID
    INTO WIN_PLAYER_USER_ID
    FROM PLAYER
    WHERE ID = WIN_PLAYER_ID;

    SELECT P.USER_ID
    INTO WIN_JUDGE_USER_ID
    FROM LOBBY AS L
        INNER JOIN JUDGE AS J ON J.LOBBY_ID = L.ID
        INNER JOIN PLAYER AS P ON P.ID = J.PLAYER_ID
    WHERE L.ID = WIN_LOBBY_ID;

    INSERT INTO WIN (LOBBY_ID, PLAYER_ID)
    VALUES (WIN_LOBBY_ID, WIN_PLAYER_ID);

    INSERT INTO LOG_WIN (PLAYER_USER_ID, JUDGE_USER_ID, CARD_ID)
    VALUES(WIN_PLAYER_USER_ID, WIN_JUDGE_USER_ID, WIN_CARD_ID);

    DELETE FROM BOARD
    WHERE LOBBY_ID = WIN_LOBBY_ID;

    SELECT U.NAME
    FROM PLAYER AS P
        INNER JOIN USER AS U ON U.ID = P.USER_ID
    WHERE P.ID = WIN_PLAYER_ID;
END;
//
DELIMITER ;

CREATE TABLE LOG_DRAW
(
    ID CHAR(36) NOT NULL DEFAULT UUID(),
    CREATED_ON_DATE DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(),

    PLAYER_USER_ID CHAR(36) NOT NULL,
    CARD_ID CHAR(36) NOT NULL,

    PRIMARY KEY (ID),
    FOREIGN KEY (PLAYER_USER_ID) REFERENCES USER (ID) ON DELETE CASCADE,
    FOREIGN KEY (CARD_ID) REFERENCES CARD (ID) ON DELETE CASCADE
);

CREATE TABLE LOG_DISCARD
(
    ID CHAR(36) NOT NULL DEFAULT UUID(),
    CREATED_ON_DATE DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(),

    PLAYER_USER_ID CHAR(36) NOT NULL,
    CARD_ID CHAR(36) NOT NULL,

    PRIMARY KEY (ID),
    FOREIGN KEY (PLAYER_USER_ID) REFERENCES USER (ID) ON DELETE CASCADE,
    FOREIGN KEY (CARD_ID) REFERENCES CARD (ID) ON DELETE CASCADE
);

CREATE TABLE LOG_PLAY
(
    ID CHAR(36) NOT NULL DEFAULT UUID(),
    CREATED_ON_DATE DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(),

    PLAYER_USER_ID CHAR(36) NOT NULL,
    JUDGE_USER_ID CHAR(36) NOT NULL,
    CARD_ID CHAR(36) NOT NULL,

    PRIMARY KEY (ID),
    FOREIGN KEY (PLAYER_USER_ID) REFERENCES USER (ID) ON DELETE CASCADE,
    FOREIGN KEY (JUDGE_USER_ID) REFERENCES USER (ID) ON DELETE CASCADE,
    FOREIGN KEY (CARD_ID) REFERENCES CARD (ID) ON DELETE CASCADE
);

CREATE TABLE LOG_WIN
(
    ID CHAR(36) NOT NULL DEFAULT UUID(),
    CREATED_ON_DATE DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(),

    PLAYER_USER_ID CHAR(36) NOT NULL,
    JUDGE_USER_ID CHAR(36) NOT NULL,
    CARD_ID CHAR(36) NOT NULL,

    PRIMARY KEY (ID),
    FOREIGN KEY (PLAYER_USER_ID) REFERENCES USER (ID) ON DELETE CASCADE,
    FOREIGN KEY (JUDGE_USER_ID) REFERENCES USER (ID) ON DELETE CASCADE,
    FOREIGN KEY (CARD_ID) REFERENCES CARD (ID) ON DELETE CASCADE
);
