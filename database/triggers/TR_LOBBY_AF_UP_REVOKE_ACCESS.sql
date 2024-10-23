CREATE TRIGGER IF NOT EXISTS TR_LOBBY_AF_UP_REVOKE_ACCESS
AFTER UPDATE ON LOBBY
FOR EACH ROW
BEGIN
    IF OLD.PASSWORD_HASH <> NEW.PASSWORD_HASH THEN
        DELETE FROM USER_ACCESS_LOBBY
        WHERE LOBBY_ID = NEW.ID;
    END IF;
END;