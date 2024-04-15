DO
$$
    BEGIN
        IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_anime_list') THEN
            CREATE TYPE user_anime_list AS ENUM (
                'WATCHING',
                'WATCHED',
                'WANT_TO_WATCH',
                'DROPPED',
                'DEFERRED',
                'REVISING'
                );
        END IF;
    END
$$;