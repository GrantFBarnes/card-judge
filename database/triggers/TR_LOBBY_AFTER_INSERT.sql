CREATE TRIGGER IF NOT EXISTS TR_LOBBY_AFTER_INSERT
    AFTER INSERT
    ON LOBBY
    FOR EACH ROW
BEGIN
    INSERT
    INTO JUDGE
        (
            LOBBY_ID
        )
    VALUES
        (
            NEW.ID
        );

    INSERT
    INTO DECK
        (
            ID,
            NAME,
            PASSWORD_HASH,
            IS_LOBBY_WILD_DECK
        )
    VALUES
        (
            NEW.ID,
            NEW.ID,
            '000000000000000000000000000000000000000000000000000000000000',
            TRUE
        );
END;