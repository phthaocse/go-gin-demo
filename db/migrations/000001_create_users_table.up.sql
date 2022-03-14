CREATE TABLE IF NOT EXISTS "user"(
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    username VARCHAR(50) NOT NULL UNIQUE,
    password TEXT NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    role VARCHAR(10),
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP
);

INSERT INTO "user"
(email, username, password, role)
VALUES ('admin@email.com', 'admin', '$2a$10$fPRhQRIV.q8h2v.jjKJPWeuapcowrI1fUtQ6VvnN35uC3PZaukLn.', 'admin');
