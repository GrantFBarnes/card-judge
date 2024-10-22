USE CARD_JUDGE;

INSERT INTO USER (NAME, PASSWORD_HASH)
VALUES ('Test User A', '$2a$14$t7gWxR3Ak8uBkyPnw4TZz.WcN3nVlbDMEQgqHOuxEfWN3yCL3dgY.');

INSERT INTO USER (NAME, PASSWORD_HASH)
VALUES ('Test User B', '$2a$14$t7gWxR3Ak8uBkyPnw4TZz.WcN3nVlbDMEQgqHOuxEfWN3yCL3dgY.');

INSERT INTO DECK (NAME, PASSWORD_HASH)
VALUES ('Test Deck A', '$2a$14$t7gWxR3Ak8uBkyPnw4TZz.WcN3nVlbDMEQgqHOuxEfWN3yCL3dgY.');

SELECT ID FROM DECK WHERE NAME = 'Test Deck A' INTO @TEST_DECK_A_ID;

INSERT INTO CARD (DECK_ID, CATEGORY, TEXT)
VALUES (@TEST_DECK_A_ID, 'PROMPT', 'Prompt Test A1'),
       (@TEST_DECK_A_ID, 'PROMPT', 'Prompt Test A2'),
       (@TEST_DECK_A_ID, 'PROMPT', 'Prompt Test A3'),
       (@TEST_DECK_A_ID, 'PROMPT', 'Prompt Test A4'),
       (@TEST_DECK_A_ID, 'PROMPT', 'Prompt Test A5'),
       (@TEST_DECK_A_ID, 'RESPONSE', 'Response Test A1'),
       (@TEST_DECK_A_ID, 'RESPONSE', 'Response Test A2'),
       (@TEST_DECK_A_ID, 'RESPONSE', 'Response Test A3'),
       (@TEST_DECK_A_ID, 'RESPONSE', 'Response Test A4'),
       (@TEST_DECK_A_ID, 'RESPONSE', 'Response Test A5'),
       (@TEST_DECK_A_ID, 'RESPONSE', 'Response Test A6'),
       (@TEST_DECK_A_ID, 'RESPONSE', 'Response Test A7'),
       (@TEST_DECK_A_ID, 'RESPONSE', 'Response Test A8'),
       (@TEST_DECK_A_ID, 'RESPONSE', 'Response Test A9'),
       (@TEST_DECK_A_ID, 'RESPONSE', 'Response Test A10'),
       (@TEST_DECK_A_ID, 'RESPONSE', 'Response Test A11'),
       (@TEST_DECK_A_ID, 'RESPONSE', 'Response Test A12');

INSERT INTO DECK (NAME, PASSWORD_HASH)
VALUES ('Test Deck B', '$2a$14$t7gWxR3Ak8uBkyPnw4TZz.WcN3nVlbDMEQgqHOuxEfWN3yCL3dgY.');

SELECT ID FROM DECK WHERE NAME = 'Test Deck B' INTO @TEST_DECK_B_ID;

INSERT INTO CARD (DECK_ID, CATEGORY, TEXT)
VALUES (@TEST_DECK_B_ID, 'PROMPT', 'Prompt Test B1'),
       (@TEST_DECK_B_ID, 'PROMPT', 'Prompt Test B2'),
       (@TEST_DECK_B_ID, 'PROMPT', 'Prompt Test B3'),
       (@TEST_DECK_B_ID, 'PROMPT', 'Prompt Test B4'),
       (@TEST_DECK_B_ID, 'PROMPT', 'Prompt Test B5'),
       (@TEST_DECK_B_ID, 'RESPONSE', 'Response Test B1'),
       (@TEST_DECK_B_ID, 'RESPONSE', 'Response Test B2'),
       (@TEST_DECK_B_ID, 'RESPONSE', 'Response Test B3'),
       (@TEST_DECK_B_ID, 'RESPONSE', 'Response Test B4'),
       (@TEST_DECK_B_ID, 'RESPONSE', 'Response Test B5'),
       (@TEST_DECK_B_ID, 'RESPONSE', 'Response Test B6'),
       (@TEST_DECK_B_ID, 'RESPONSE', 'Response Test B7'),
       (@TEST_DECK_B_ID, 'RESPONSE', 'Response Test B8'),
       (@TEST_DECK_B_ID, 'RESPONSE', 'Response Test B9'),
       (@TEST_DECK_B_ID, 'RESPONSE', 'Response Test B10'),
       (@TEST_DECK_B_ID, 'RESPONSE', 'Response Test B11'),
       (@TEST_DECK_B_ID, 'RESPONSE', 'Response Test B12');
