USE CARD_JUDGE;

INSERT INTO DECK (ID, NAME)
VALUES (UUID(), 'Test 1'),
       (UUID(), 'Test 2');


SELECT ID FROM DECK WHERE NAME = 'Test 1' into @DeckId1;
SELECT ID FROM CARD_TYPE WHERE NAME = 'Judge' into @JudgeCardId;
SELECT ID FROM CARD_TYPE WHERE NAME = 'Player' into @PlayerCardId;


INSERT INTO CARD (ID, DECK_ID, CARD_TYPE_ID, TEXT)
VALUES (UUID(), @DeckId1, @JudgeCardId,  'Judge Test 1'),
       (UUID(), @DeckId1, @JudgeCardId,  'Judge Test 2'),
       (UUID(), @DeckId1, @JudgeCardId,  'Judge Test 3'),
       (UUID(), @DeckId1, @JudgeCardId,  'Judge Test 4'),
       (UUID(), @DeckId1, @JudgeCardId,  'Judge Test 5'),
       (UUID(), @DeckId1, @JudgeCardId,  'Judge Test 6'),
       (UUID(), @DeckId1, @JudgeCardId,  'Judge Test 8'),
       (UUID(), @DeckId1, @JudgeCardId,  'Judge Test 9'),
       (UUID(), @DeckId1, @JudgeCardId,  'Judge Test 10'),
       (UUID(), @DeckId1, @JudgeCardId,  'Judge Test 11'),
       (UUID(), @DeckId1, @JudgeCardId,  'Judge Test 12'),
       (UUID(), @DeckId1, @PlayerCardId,  'Player Test 1'),
       (UUID(), @DeckId1, @PlayerCardId,  'Player Test 2'),
       (UUID(), @DeckId1, @PlayerCardId,  'Player Test 3'),
       (UUID(), @DeckId1, @PlayerCardId,  'Player Test 4'),
       (UUID(), @DeckId1, @PlayerCardId,  'Player Test 5'),
       (UUID(), @DeckId1, @PlayerCardId,  'Player Test 6'),
       (UUID(), @DeckId1, @PlayerCardId,  'Player Test 8'),
       (UUID(), @DeckId1, @PlayerCardId,  'Player Test 9'),
       (UUID(), @DeckId1, @PlayerCardId,  'Player Test 10'),
       (UUID(), @DeckId1, @PlayerCardId,  'Player Test 11'),
       (UUID(), @DeckId1, @PlayerCardId,  'Player Test 12');





















