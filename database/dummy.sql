INSERT INTO DECK (ID, NAME)
VALUES (uuid(), 'Test 1'),
       (uuid(), 'Test 2');


select ID FROM DECK WHERE NAME = 'Test 1' into @DeckId1;
SELECT ID FROM CARD_TYPE WHERE NAME = 'Judge' into @JudgeCardId;
SELECT ID FROM CARD_TYPE WHERE NAME = 'Player' into @PlayerCardId;


INSERT INTO CARD (ID, DECK_ID, CARD_TYPE_ID, TEXT)
VALUES (uuid(), @DeckId1, @JudgeCardId,  'Judge Test 1'),
       (uuid(), @DeckId1, @JudgeCardId,  'Judge Test 2'),
       (uuid(), @DeckId1, @JudgeCardId,  'Judge Test 3'),
       (uuid(), @DeckId1, @JudgeCardId,  'Judge Test 4'),
       (uuid(), @DeckId1, @JudgeCardId,  'Judge Test 5'),
       (uuid(), @DeckId1, @JudgeCardId,  'Judge Test 6'),
       (uuid(), @DeckId1, @JudgeCardId,  'Judge Test 8'),
       (uuid(), @DeckId1, @JudgeCardId,  'Judge Test 9'),
       (uuid(), @DeckId1, @JudgeCardId,  'Judge Test 10'),
       (uuid(), @DeckId1, @JudgeCardId,  'Judge Test 11'),
       (uuid(), @DeckId1, @JudgeCardId,  'Judge Test 12'),
       (uuid(), @DeckId1, @PlayerCardId,  'Player Test 1'),
       (uuid(), @DeckId1, @PlayerCardId,  'Player Test 2'),
       (uuid(), @DeckId1, @PlayerCardId,  'Player Test 3'),
       (uuid(), @DeckId1, @PlayerCardId,  'Player Test 4'),
       (uuid(), @DeckId1, @PlayerCardId,  'Player Test 5'),
       (uuid(), @DeckId1, @PlayerCardId,  'Player Test 6'),
       (uuid(), @DeckId1, @PlayerCardId,  'Player Test 8'),
       (uuid(), @DeckId1, @PlayerCardId,  'Player Test 9'),
       (uuid(), @DeckId1, @PlayerCardId,  'Player Test 10'),
       (uuid(), @DeckId1, @PlayerCardId,  'Player Test 11'),
       (uuid(), @DeckId1, @PlayerCardId,  'Player Test 12');





















