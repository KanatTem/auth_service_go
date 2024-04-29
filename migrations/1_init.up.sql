CREATE TABLE IF NOT EXISTS users (
                                     id SERIAL PRIMARY KEY,
                                     email TEXT NOT NULL UNIQUE,
                                     pass_hash BYTEA NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_email ON users (email);

CREATE TABLE IF NOT EXISTS apps (
                                    id SERIAL PRIMARY KEY,
                                    name TEXT NOT NULL UNIQUE,
                                    secret TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS roles (
                                     id   SERIAL    PRIMARY KEY,
                                     app_id INT NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
                                     name TEXT NOT NULL,
                                     UNIQUE (app_id, name)
);


CREATE TABLE IF NOT EXISTS user_roles (
                                          user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                          role_id INT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
                                          PRIMARY KEY (user_id, role_id)
);
