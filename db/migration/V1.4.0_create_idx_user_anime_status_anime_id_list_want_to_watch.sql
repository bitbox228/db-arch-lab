CREATE INDEX idx_user_anime_status_anime_id_list_partial ON user_anime_status (anime_id) WHERE list <> 'WANT_TO_WATCH';