CREATE TRIGGER IF NOT EXISTS TR_AUDIT_CARD_UPDATE
BEFORE UPDATE ON CARD
FOR EACH ROW
BEGIN
    INSERT INTO AUDIT_CARD(
        AUDIT_TYPE,
        CARD_ID,
        DECK_ID,
        CATEGORY,
        TEXT,
        YOUTUBE,
        IMAGE
    )
    VALUES (
        'UPDATE',
        OLD.ID,
        OLD.DECK_ID,
        OLD.CATEGORY,
        OLD.TEXT,
        OLD.YOUTUBE,
        OLD.IMAGE
    );
END;