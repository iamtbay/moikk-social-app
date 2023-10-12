CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users(
    "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
    "name" VARCHAR(150) NOT NULL,
    "username" VARCHAR(150) NOT NULL UNIQUE,
    "email" VARCHAR(150) NOT NULL UNIQUE,
    "password" VARCHAR NOT NULL,
    "location" VARCHAR(150),
    "profile_photo" VARCHAR DEFAULT 'nophoto'
);