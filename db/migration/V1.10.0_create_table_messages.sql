CREATE TABLE IF NOT EXISTS messages
(
    id          SERIAL PRIMARY KEY,
    sender_id   INTEGER REFERENCES users (id) NOT NULL,
    receiver_id INTEGER REFERENCES users (id) NOT NULL,
    text        VARCHAR(1000)                 NOT NULL,
    file_url    VARCHAR(255)                  NOT NULL,
    time        TIMESTAMP                     NOT NULL
);