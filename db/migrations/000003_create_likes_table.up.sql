CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS likes(
    id uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
    post_id uuid NOT NULL,
    user_id uuid NOT NULL,
    created_at timestamp with time zone default now(),
    CONSTRAINT fk_postid
        FOREIGN KEY (post_id)
            REFERENCES posts(id)
                ON DELETE CASCADE,
    CONSTRAINT fk_userid
        FOREIGN KEY (user_id)
            REFERENCES users(id)
                ON DELETE CASCADE
);