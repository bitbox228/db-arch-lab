BEGIN;
CREATE TABLE IF NOT EXISTS reactions_partition
(
    id        SERIAL,
    review_id INTEGER REFERENCES reviews (id) NOT NULL,
    user_id   INTEGER REFERENCES users (id)   NOT NULL,
    is_like   BOOLEAN                         NOT NULL,
    UNIQUE (id, is_like)
) PARTITION BY LIST (is_like);

CREATE TABLE IF NOT EXISTS reactions_like PARTITION OF reactions_partition
    FOR VALUES IN (TRUE);

CREATE TABLE IF NOT EXISTS reactions_dislike PARTITION OF reactions_partition
    FOR VALUES IN (FALSE);

DO
$$
    BEGIN
        IF NOT EXISTS (SELECT 1 FROM reactions_partition) THEN
            INSERT INTO reactions_partition (id, review_id, user_id, is_like)
            SELECT id, review_id, user_id, is_like
            FROM reactions;
        END IF;
    END
$$;

CREATE INDEX IF NOT EXISTS idx_reactions_dislike_review_id ON reactions_dislike (review_id);
CREATE INDEX IF NOT EXISTS idx_reactions_dislike_is_like ON reactions_dislike (is_like);

CREATE INDEX IF NOT EXISTS idx_reactions_like_review_id ON reactions_like (review_id);
CREATE INDEX IF NOT EXISTS idx_reactions_like_is_like ON reactions_like (is_like);

DROP TABLE IF EXISTS reactions;
ALTER TABLE IF EXISTS reactions_partition
    RENAME TO reactions;

END;