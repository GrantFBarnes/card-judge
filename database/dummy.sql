USE CARD_JUDGE;

INSERT INTO USER (NAME, PASSWORD_HASH)
VALUES ('Test User A', '$2a$14$t7gWxR3Ak8uBkyPnw4TZz.WcN3nVlbDMEQgqHOuxEfWN3yCL3dgY.');

INSERT INTO USER (NAME, PASSWORD_HASH)
VALUES ('Test User B', '$2a$14$t7gWxR3Ak8uBkyPnw4TZz.WcN3nVlbDMEQgqHOuxEfWN3yCL3dgY.');

INSERT INTO DECK (NAME, PASSWORD_HASH)
VALUES ('Test Deck A', '$2a$14$t7gWxR3Ak8uBkyPnw4TZz.WcN3nVlbDMEQgqHOuxEfWN3yCL3dgY.');

SELECT ID FROM DECK WHERE NAME = 'Test Deck A' INTO @TEST_DECK_A_ID;

INSERT INTO CARD (DECK_ID, CATEGORY, TEXT)
VALUES (@TEST_DECK_A_ID, 'JUDGE', 'Judge Test A1'),
       (@TEST_DECK_A_ID, 'JUDGE', 'Judge Test A2'),
       (@TEST_DECK_A_ID, 'JUDGE', 'Judge Test A3'),
       (@TEST_DECK_A_ID, 'JUDGE', 'Judge Test A4'),
       (@TEST_DECK_A_ID, 'JUDGE', 'Judge Test A5'),
       (@TEST_DECK_A_ID, 'PLAYER', 'Player Test A1'),
       (@TEST_DECK_A_ID, 'PLAYER', 'Player Test A2'),
       (@TEST_DECK_A_ID, 'PLAYER', 'Player Test A3'),
       (@TEST_DECK_A_ID, 'PLAYER', 'Player Test A4'),
       (@TEST_DECK_A_ID, 'PLAYER', 'Player Test A5'),
       (@TEST_DECK_A_ID, 'PLAYER', 'Player Test A6'),
       (@TEST_DECK_A_ID, 'PLAYER', 'Player Test A7'),
       (@TEST_DECK_A_ID, 'PLAYER', 'Player Test A8'),
       (@TEST_DECK_A_ID, 'PLAYER', 'Player Test A9'),
       (@TEST_DECK_A_ID, 'PLAYER', 'Player Test A10'),
       (@TEST_DECK_A_ID, 'PLAYER', 'Player Test A11'),
       (@TEST_DECK_A_ID, 'PLAYER', 'Player Test A12');

INSERT INTO DECK (NAME, PASSWORD_HASH)
VALUES ('Test Deck B', '$2a$14$t7gWxR3Ak8uBkyPnw4TZz.WcN3nVlbDMEQgqHOuxEfWN3yCL3dgY.');

SELECT ID FROM DECK WHERE NAME = 'Test Deck B' INTO @TEST_DECK_B_ID;

INSERT INTO CARD (DECK_ID, CATEGORY, TEXT)
VALUES (@TEST_DECK_B_ID, 'JUDGE', 'Judge Test B1'),
       (@TEST_DECK_B_ID, 'JUDGE', 'Judge Test B2'),
       (@TEST_DECK_B_ID, 'JUDGE', 'Judge Test B3'),
       (@TEST_DECK_B_ID, 'JUDGE', 'Judge Test B4'),
       (@TEST_DECK_B_ID, 'JUDGE', 'Judge Test B5'),
       (@TEST_DECK_B_ID, 'PLAYER', 'Player Test B1'),
       (@TEST_DECK_B_ID, 'PLAYER', 'Player Test B2'),
       (@TEST_DECK_B_ID, 'PLAYER', 'Player Test B3'),
       (@TEST_DECK_B_ID, 'PLAYER', 'Player Test B4'),
       (@TEST_DECK_B_ID, 'PLAYER', 'Player Test B5'),
       (@TEST_DECK_B_ID, 'PLAYER', 'Player Test B6'),
       (@TEST_DECK_B_ID, 'PLAYER', 'Player Test B7'),
       (@TEST_DECK_B_ID, 'PLAYER', 'Player Test B8'),
       (@TEST_DECK_B_ID, 'PLAYER', 'Player Test B9'),
       (@TEST_DECK_B_ID, 'PLAYER', 'Player Test B10'),
       (@TEST_DECK_B_ID, 'PLAYER', 'Player Test B11'),
       (@TEST_DECK_B_ID, 'PLAYER', 'Player Test B12');
