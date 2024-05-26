DO
$$
    BEGIN
        IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'anime_status') THEN
            CREATE TYPE anime_status AS ENUM (
                'ONGOING',
                'PLANNED',
                'RELEASED'
                );
        END IF;
    END
$$;