CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS posts(
    id uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
    user_id uuid NOT NULL,
    title   VARCHAR NOT NULL,
    content VARCHAR NOT NULL,
    files varchar[],
    created_at timestamp with time zone DEFAULT Now() not null,
    updated_at timestamp with time zone,
    CONSTRAINT fk_userid
        FOREIGN KEY (user_id)
            REFERENCES users(id)
                ON DELETE CASCADE
);