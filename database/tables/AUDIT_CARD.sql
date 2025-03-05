CREATE TABLE IF NOT EXISTS AUDIT_CARD
(
    ID              UUID                       NOT NULL DEFAULT UUID(),
    CREATED_ON_DATE DATETIME                   NOT NULL DEFAULT CURRENT_TIMESTAMP(),

    AUDIT_TYPE      ENUM ('UPDATE','DELETE')   NOT NULL,
    CARD_ID         UUID                       NOT NULL,
    DECK_ID         UUID                       NOT NULL,
    CATEGORY        ENUM ('PROMPT','RESPONSE') NOT NULL,
    TEXT            VARCHAR(510)               NOT NULL,
    IMAGE           BLOB                       NULL,

    PRIMARY KEY (ID)
);