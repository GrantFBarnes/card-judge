CREATE TRIGGER IF NOT EXISTS TR_AUDIT_CARD_DELETE
    BEFORE DELETE
    ON CARD
    FOR EACH ROW
BEGIN
    INSERT INTO AUDIT_CARD (
        AUDIT_TYPE,
        CARD_ID,
        DECK_ID,
        CATEGORY,
        TEXT,
        IMAGE
    ) VALUES (
        'DELETE',
        OLD.ID,
        OLD.DECK_ID,
        OLD.CATEGORY,
        OLD.TEXT,
        OLD.IMAGE
    );
END;
