CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS follows(
    id uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
    follower_user_id uuid NOT NULL,
    followed_user_id uuid NOT NULL,
    created_at timestamp with time zone default now(),
    CONSTRAINT fk_follower
        FOREIGN KEY(follower_user_id)
            REFERENCES users(id)
                ON DELETE CASCADE,
    CONSTRAINT fk_followed
        FOREIGN KEY(followed_user_id)
            REFERENCES users(id)
                ON DELETE CASCADE
);

