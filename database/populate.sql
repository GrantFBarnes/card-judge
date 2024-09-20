USE CARD_JUDGE;

INSERT INTO USER (NAME, PASSWORD_HASH, IS_ADMIN)
VALUES ('Grant', '$2a$14$t7gWxR3Ak8uBkyPnw4TZz.WcN3nVlbDMEQgqHOuxEfWN3yCL3dgY.', 1);

INSERT INTO CARD_TYPE (ID, NAME)
VALUES (uuid(), 'Judge'),
       (uuid(), 'Player');




