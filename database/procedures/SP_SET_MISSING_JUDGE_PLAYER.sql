CREATE PROCEDURE IF NOT EXISTS SP_SET_MISSING_JUDGE_PLAYER(IN VAR_LOBBY_ID UUID)
BEGIN
    IF FN_GET_LOBBY_JUDGE_PLAYER_ID(VAR_LOBBY_ID) IS NULL THEN
        CALL SP_SET_NEXT_JUDGE_PLAYER(VAR_LOBBY_ID);
    END
    IF;
END;