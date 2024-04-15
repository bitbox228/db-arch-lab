CREATE TABLE IF NOT EXISTS anime_series
(
    id            SERIAL PRIMARY KEY,
    anime_id      INTEGER REFERENCES anime (id) NOT NULL,
    series_url    VARCHAR(255)                  NOT NULL,
    seconds_count INTEGER                       NOT NULL
);