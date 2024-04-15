DO
$$
    BEGIN
        IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'notification_type') THEN
            CREATE TYPE notification_type AS ENUM (
                'FRIEND_REQUEST',
                'NEW_EPISODE',
                'NEW_MESSAGE'
                );
        END IF;
    END
$$;