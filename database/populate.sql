USE CARD_JUDGE;

INSERT INTO PLAYER (NAME, PASSWORD_HASH, IS_ADMIN)
VALUES ('Grant', '$2a$14$t7gWxR3Ak8uBkyPnw4TZz.WcN3nVlbDMEQgqHOuxEfWN3yCL3dgY.', 1);

INSERT INTO DECK (ID, NAME, PASSWORD_HASH)
VALUES ('f395b797-6d89-11ef-aad4-28800dbd8d8a', 'Deck One', NULL),
       ('f395b862-6d89-11ef-aad4-28800dbd8d8a', 'Deck Two', NULL),
       ('f395b8e4-6d89-11ef-aad4-28800dbd8d8a', 'Deck Three', '$2a$14$t7gWxR3Ak8uBkyPnw4TZz.WcN3nVlbDMEQgqHOuxEfWN3yCL3dgY.');

INSERT INTO CARD_TYPE (ID, NAME)
VALUES ('a907026b-6fa6-11ef-b1ac-3bd680fc6f38', 'Judge'),
       ('a90703f8-6fa6-11ef-b1ac-3bd680fc6f38', 'Player');

INSERT INTO CARD (DECK_ID, CARD_TYPE_ID, TEXT)
VALUES ('f395b797-6d89-11ef-aad4-28800dbd8d8a', 'a907026b-6fa6-11ef-b1ac-3bd680fc6f38', 'Deck One - Judge Card 1'),
       ('f395b797-6d89-11ef-aad4-28800dbd8d8a', 'a907026b-6fa6-11ef-b1ac-3bd680fc6f38', 'Deck One - Judge Card 2'),
       ('f395b797-6d89-11ef-aad4-28800dbd8d8a', 'a907026b-6fa6-11ef-b1ac-3bd680fc6f38', 'Deck One - Judge Card 3'),
       ('f395b797-6d89-11ef-aad4-28800dbd8d8a', 'a907026b-6fa6-11ef-b1ac-3bd680fc6f38', 'Deck One - Judge Card 4'),
       ('f395b797-6d89-11ef-aad4-28800dbd8d8a', 'a907026b-6fa6-11ef-b1ac-3bd680fc6f38', 'Deck One - Judge Card 5'),
       ('f395b797-6d89-11ef-aad4-28800dbd8d8a', 'a90703f8-6fa6-11ef-b1ac-3bd680fc6f38', 'Deck One - Player Card 1'),
       ('f395b797-6d89-11ef-aad4-28800dbd8d8a', 'a90703f8-6fa6-11ef-b1ac-3bd680fc6f38', 'Deck One - Player Card 2'),
       ('f395b797-6d89-11ef-aad4-28800dbd8d8a', 'a90703f8-6fa6-11ef-b1ac-3bd680fc6f38', 'Deck One - Player Card 3'),
       ('f395b797-6d89-11ef-aad4-28800dbd8d8a', 'a90703f8-6fa6-11ef-b1ac-3bd680fc6f38', 'Deck One - Player Card 4'),
       ('f395b797-6d89-11ef-aad4-28800dbd8d8a', 'a90703f8-6fa6-11ef-b1ac-3bd680fc6f38', 'Deck One - Player Card 5'),
       ('f395b797-6d89-11ef-aad4-28800dbd8d8a', 'a90703f8-6fa6-11ef-b1ac-3bd680fc6f38', 'Deck One - Player Card 6'),
       ('f395b797-6d89-11ef-aad4-28800dbd8d8a', 'a90703f8-6fa6-11ef-b1ac-3bd680fc6f38', 'Deck One - Player Card 7'),
       ('f395b797-6d89-11ef-aad4-28800dbd8d8a', 'a90703f8-6fa6-11ef-b1ac-3bd680fc6f38', 'Deck One - Player Card 8'),
       ('f395b797-6d89-11ef-aad4-28800dbd8d8a', 'a90703f8-6fa6-11ef-b1ac-3bd680fc6f38', 'Deck One - Player Card 9'),
       ('f395b797-6d89-11ef-aad4-28800dbd8d8a', 'a90703f8-6fa6-11ef-b1ac-3bd680fc6f38', 'Deck One - Player Card 10');

INSERT INTO CARD (DECK_ID, CARD_TYPE_ID, TEXT)
VALUES ('f395b862-6d89-11ef-aad4-28800dbd8d8a', 'a907026b-6fa6-11ef-b1ac-3bd680fc6f38', 'Deck Two - Judge Card 1'),
       ('f395b862-6d89-11ef-aad4-28800dbd8d8a', 'a907026b-6fa6-11ef-b1ac-3bd680fc6f38', 'Deck Two - Judge Card 2'),
       ('f395b862-6d89-11ef-aad4-28800dbd8d8a', 'a907026b-6fa6-11ef-b1ac-3bd680fc6f38', 'Deck Two - Judge Card 3'),
       ('f395b862-6d89-11ef-aad4-28800dbd8d8a', 'a907026b-6fa6-11ef-b1ac-3bd680fc6f38', 'Deck Two - Judge Card 4'),
       ('f395b862-6d89-11ef-aad4-28800dbd8d8a', 'a907026b-6fa6-11ef-b1ac-3bd680fc6f38', 'Deck Two - Judge Card 5'),
       ('f395b862-6d89-11ef-aad4-28800dbd8d8a', 'a90703f8-6fa6-11ef-b1ac-3bd680fc6f38', 'Deck Two - Player Card 1'),
       ('f395b862-6d89-11ef-aad4-28800dbd8d8a', 'a90703f8-6fa6-11ef-b1ac-3bd680fc6f38', 'Deck Two - Player Card 2'),
       ('f395b862-6d89-11ef-aad4-28800dbd8d8a', 'a90703f8-6fa6-11ef-b1ac-3bd680fc6f38', 'Deck Two - Player Card 3'),
       ('f395b862-6d89-11ef-aad4-28800dbd8d8a', 'a90703f8-6fa6-11ef-b1ac-3bd680fc6f38', 'Deck Two - Player Card 4'),
       ('f395b862-6d89-11ef-aad4-28800dbd8d8a', 'a90703f8-6fa6-11ef-b1ac-3bd680fc6f38', 'Deck Two - Player Card 5'),
       ('f395b862-6d89-11ef-aad4-28800dbd8d8a', 'a90703f8-6fa6-11ef-b1ac-3bd680fc6f38', 'Deck Two - Player Card 6'),
       ('f395b862-6d89-11ef-aad4-28800dbd8d8a', 'a90703f8-6fa6-11ef-b1ac-3bd680fc6f38', 'Deck Two - Player Card 7'),
       ('f395b862-6d89-11ef-aad4-28800dbd8d8a', 'a90703f8-6fa6-11ef-b1ac-3bd680fc6f38', 'Deck Two - Player Card 8'),
       ('f395b862-6d89-11ef-aad4-28800dbd8d8a', 'a90703f8-6fa6-11ef-b1ac-3bd680fc6f38', 'Deck Two - Player Card 9'),
       ('f395b862-6d89-11ef-aad4-28800dbd8d8a', 'a90703f8-6fa6-11ef-b1ac-3bd680fc6f38', 'Deck Two - Player Card 10');

INSERT INTO CARD (DECK_ID, CARD_TYPE_ID, TEXT)
VALUES ('f395b8e4-6d89-11ef-aad4-28800dbd8d8a', 'a907026b-6fa6-11ef-b1ac-3bd680fc6f38', 'Deck Three - Judge Card 1'),
       ('f395b8e4-6d89-11ef-aad4-28800dbd8d8a', 'a907026b-6fa6-11ef-b1ac-3bd680fc6f38', 'Deck Three - Judge Card 2'),
       ('f395b8e4-6d89-11ef-aad4-28800dbd8d8a', 'a907026b-6fa6-11ef-b1ac-3bd680fc6f38', 'Deck Three - Judge Card 3'),
       ('f395b8e4-6d89-11ef-aad4-28800dbd8d8a', 'a907026b-6fa6-11ef-b1ac-3bd680fc6f38', 'Deck Three - Judge Card 4'),
       ('f395b8e4-6d89-11ef-aad4-28800dbd8d8a', 'a907026b-6fa6-11ef-b1ac-3bd680fc6f38', 'Deck Three - Judge Card 5'),
       ('f395b8e4-6d89-11ef-aad4-28800dbd8d8a', 'a90703f8-6fa6-11ef-b1ac-3bd680fc6f38', 'Deck Three - Player Card 1'),
       ('f395b8e4-6d89-11ef-aad4-28800dbd8d8a', 'a90703f8-6fa6-11ef-b1ac-3bd680fc6f38', 'Deck Three - Player Card 2'),
       ('f395b8e4-6d89-11ef-aad4-28800dbd8d8a', 'a90703f8-6fa6-11ef-b1ac-3bd680fc6f38', 'Deck Three - Player Card 3'),
       ('f395b8e4-6d89-11ef-aad4-28800dbd8d8a', 'a90703f8-6fa6-11ef-b1ac-3bd680fc6f38', 'Deck Three - Player Card 4'),
       ('f395b8e4-6d89-11ef-aad4-28800dbd8d8a', 'a90703f8-6fa6-11ef-b1ac-3bd680fc6f38', 'Deck Three - Player Card 5'),
       ('f395b8e4-6d89-11ef-aad4-28800dbd8d8a', 'a90703f8-6fa6-11ef-b1ac-3bd680fc6f38', 'Deck Three - Player Card 6'),
       ('f395b8e4-6d89-11ef-aad4-28800dbd8d8a', 'a90703f8-6fa6-11ef-b1ac-3bd680fc6f38', 'Deck Three - Player Card 7'),
       ('f395b8e4-6d89-11ef-aad4-28800dbd8d8a', 'a90703f8-6fa6-11ef-b1ac-3bd680fc6f38', 'Deck Three - Player Card 8'),
       ('f395b8e4-6d89-11ef-aad4-28800dbd8d8a', 'a90703f8-6fa6-11ef-b1ac-3bd680fc6f38', 'Deck Three - Player Card 9'),
       ('f395b8e4-6d89-11ef-aad4-28800dbd8d8a', 'a90703f8-6fa6-11ef-b1ac-3bd680fc6f38', 'Deck Three - Player Card 10');
