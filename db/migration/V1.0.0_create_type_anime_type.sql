DO
$$
    BEGIN
        IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'anime_type') THEN
            CREATE TYPE anime_type AS ENUM (
                'SERIES',
                'MOVIE',
                'OVA'
                );
        END IF;
    END
$$;