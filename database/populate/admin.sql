IF NOT EXISTS(
    SELECT ID
    FROM USER
    WHERE IS_ADMIN = 1
) THEN
    -- CREATE ADMIN ACCOUNT WITH DEFAULT PASSWORD
    INSERT INTO USER (NAME, PASSWORD_HASH, IS_ADMIN)
    VALUES ('Admin', '$2a$14$t7gWxR3Ak8uBkyPnw4TZz.WcN3nVlbDMEQgqHOuxEfWN3yCL3dgY.', 1);
END IF;