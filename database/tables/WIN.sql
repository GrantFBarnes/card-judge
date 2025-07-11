CREATE TABLE IF NOT EXISTS WIN(
    ID UUID NOT NULL DEFAULT UUID(),
    CREATED_ON_DATE DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(),
    PLAYER_ID UUID NOT NULL,
    PRIMARY KEY(ID),
    FOREIGN KEY(PLAYER_ID) REFERENCES PLAYER(ID) ON DELETE CASCADE
);