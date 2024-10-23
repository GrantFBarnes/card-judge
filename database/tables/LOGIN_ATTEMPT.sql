CREATE TABLE IF NOT EXISTS LOGIN_ATTEMPT
(
    ID UUID NOT NULL DEFAULT UUID(),
    CREATED_ON_DATE DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(),

    IP_ADDRESS VARCHAR(255) NOT NULL,
    USER_NAME VARCHAR(255) NOT NULL,

    PRIMARY KEY (ID)
);