CREATE TABLE IF NOT EXISTS LOG_RESPONSE_CARD
(
    ID               UUID                             NOT NULL DEFAULT UUID(),
    CREATED_ON_DATE  DATETIME                         NOT NULL DEFAULT CURRENT_TIMESTAMP(),

    LOBBY_ID         UUID                             NOT NULL,
    ROUND_ID         UUID                             NOT NULL,
    RESPONSE_ID      UUID                             NOT NULL,
    RESPONSE_CARD_ID UUID                             NOT NULL,
    JUDGE_USER_ID    UUID                             NOT NULL,
    JUDGE_CARD_ID    UUID                             NOT NULL,
    PLAYER_USER_ID   UUID                             NOT NULL,
    PLAYER_CARD_ID   UUID                             NOT NULL,
    SPECIAL_CATEGORY ENUM ('STEAL','SURPRISE','WILD') NULL,

    PRIMARY KEY (ID)
);