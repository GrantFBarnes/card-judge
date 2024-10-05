USE CARD_JUDGE;

INSERT INTO USER (NAME, PASSWORD_HASH, IS_ADMIN)
VALUES ('Grant', '$2a$14$t7gWxR3Ak8uBkyPnw4TZz.WcN3nVlbDMEQgqHOuxEfWN3yCL3dgY.', 1);

INSERT INTO USER (NAME, PASSWORD_HASH, IS_ADMIN)
VALUES ('TestA', '$2a$14$t7gWxR3Ak8uBkyPnw4TZz.WcN3nVlbDMEQgqHOuxEfWN3yCL3dgY.', 0);

INSERT INTO USER (NAME, PASSWORD_HASH, IS_ADMIN)
VALUES ('TestB', '$2a$14$t7gWxR3Ak8uBkyPnw4TZz.WcN3nVlbDMEQgqHOuxEfWN3yCL3dgY.', 0);

INSERT INTO CARD_TYPE (NAME)
VALUES ('Judge'),
       ('Player');

INSERT INTO DECK (NAME)
VALUES ('Test 1'),
       ('Test 2');

SELECT ID FROM DECK WHERE NAME = 'Test 1' into @DeckId1;
SELECT ID FROM CARD_TYPE WHERE NAME = 'Judge' into @JudgeCardId;
SELECT ID FROM CARD_TYPE WHERE NAME = 'Player' into @PlayerCardId;

INSERT INTO CARD (DECK_ID, CARD_TYPE_ID, TEXT)
VALUES (@DeckId1, @JudgeCardId, 'Judge Test 1'),
       (@DeckId1, @JudgeCardId, 'Judge Test 2'),
       (@DeckId1, @JudgeCardId, 'Judge Test 3'),
       (@DeckId1, @JudgeCardId, 'Judge Test 4'),
       (@DeckId1, @JudgeCardId, 'Judge Test 5'),
       (@DeckId1, @JudgeCardId, 'Judge Test 6'),
       (@DeckId1, @JudgeCardId, 'Judge Test 8'),
       (@DeckId1, @JudgeCardId, 'Judge Test 9'),
       (@DeckId1, @JudgeCardId, 'Judge Test 10'),
       (@DeckId1, @JudgeCardId, 'Judge Test 11'),
       (@DeckId1, @JudgeCardId, 'Judge Test 12'),
       (@DeckId1, @PlayerCardId, 'Player Test 1'),
       (@DeckId1, @PlayerCardId, 'Player Test 2'),
       (@DeckId1, @PlayerCardId, 'Player Test 3'),
       (@DeckId1, @PlayerCardId, 'Player Test 4'),
       (@DeckId1, @PlayerCardId, 'Player Test 5'),
       (@DeckId1, @PlayerCardId, 'Player Test 6'),
       (@DeckId1, @PlayerCardId, 'Player Test 8'),
       (@DeckId1, @PlayerCardId, 'Player Test 9'),
       (@DeckId1, @PlayerCardId, 'Player Test 10'),
       (@DeckId1, @PlayerCardId, 'Player Test 11'),
       (@DeckId1, @PlayerCardId, 'Player Test 12');
