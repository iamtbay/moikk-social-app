CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS comments(
    id uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
    post_id uuid NOT NULL,
    upper_comment_id uuid,
    user_id uuid NOT NULL,
    content VARCHAR NOT NULL,
    created_at timestamp with time zone default now(),
    CONSTRAINT fk_postid 
        FOREIGN KEY (post_id)
            REFERENCES posts(id)
                ON DELETE CASCADE,
    CONSTRAINT fk_comment_id
        FOREIGN KEY (upper_comment_id)
            REFERENCES comments(id)
                ON DELETE CASCADE
);