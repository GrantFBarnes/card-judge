USE CARD_JUDGE;

INSERT INTO USER (NAME, PASSWORD_HASH, IS_ADMIN)
VALUES ('Grant', '$2a$14$t7gWxR3Ak8uBkyPnw4TZz.WcN3nVlbDMEQgqHOuxEfWN3yCL3dgY.', 1);

INSERT INTO USER (NAME, PASSWORD_HASH, IS_ADMIN)
VALUES ('TestA', '$2a$14$t7gWxR3Ak8uBkyPnw4TZz.WcN3nVlbDMEQgqHOuxEfWN3yCL3dgY.', 0);

INSERT INTO USER (NAME, PASSWORD_HASH, IS_ADMIN)
VALUES ('TestB', '$2a$14$t7gWxR3Ak8uBkyPnw4TZz.WcN3nVlbDMEQgqHOuxEfWN3yCL3dgY.', 0);

INSERT INTO DECK (NAME, PASSWORD_HASH)
VALUES ('Test 1', '$2a$14$t7gWxR3Ak8uBkyPnw4TZz.WcN3nVlbDMEQgqHOuxEfWN3yCL3dgY.'),
       ('Test 2', '$2a$14$t7gWxR3Ak8uBkyPnw4TZz.WcN3nVlbDMEQgqHOuxEfWN3yCL3dgY.');

SELECT ID FROM DECK WHERE NAME = 'Test 1' into @DeckId1;

INSERT INTO CARD (DECK_ID, CATEGORY, TEXT)
VALUES (@DeckId1, 'JUDGE', 'Judge Test 1'),
       (@DeckId1, 'JUDGE', 'Judge Test 2'),
       (@DeckId1, 'JUDGE', 'Judge Test 3'),
       (@DeckId1, 'JUDGE', 'Judge Test 4'),
       (@DeckId1, 'JUDGE', 'Judge Test 5'),
       (@DeckId1, 'JUDGE', 'Judge Test 6'),
       (@DeckId1, 'JUDGE', 'Judge Test 7'),
       (@DeckId1, 'JUDGE', 'Judge Test 8'),
       (@DeckId1, 'JUDGE', 'Judge Test 9'),
       (@DeckId1, 'JUDGE', 'Judge Test 10'),
       (@DeckId1, 'JUDGE', 'Judge Test 11'),
       (@DeckId1, 'JUDGE', 'Judge Test 12'),
       (@DeckId1, 'PLAYER', 'Player Test 1'),
       (@DeckId1, 'PLAYER', 'Player Test 2'),
       (@DeckId1, 'PLAYER', 'Player Test 3'),
       (@DeckId1, 'PLAYER', 'Player Test 4'),
       (@DeckId1, 'PLAYER', 'Player Test 5'),
       (@DeckId1, 'PLAYER', 'Player Test 6'),
       (@DeckId1, 'PLAYER', 'Player Test 7'),
       (@DeckId1, 'PLAYER', 'Player Test 8'),
       (@DeckId1, 'PLAYER', 'Player Test 9'),
       (@DeckId1, 'PLAYER', 'Player Test 10'),
       (@DeckId1, 'PLAYER', 'Player Test 11'),
       (@DeckId1, 'PLAYER', 'Player Test 12');
