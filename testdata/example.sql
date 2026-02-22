-- CodeOwner: @sql_owner

-- This is a single-line comment

/*
This is a multi-line
block comment in SQL.
*/

CREATE TABLE users (
    id INTEGER PRIMARY KEY,  -- Inline comment
    name TEXT NOT NULL,
    email TEXT UNIQUE
);

-- Another single-line comment
INSERT INTO users (name, email) VALUES ('Alice', 'alice@example.com');

SELECT * FROM users;  /* Inline block comment */
