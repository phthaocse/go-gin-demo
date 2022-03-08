CREATE TABLE IF NOT EXISTS "user"(
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL,
    username VARCHAR(50) NOT NULL,
    password TEXT NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    role VARCHAR(10),
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP
);