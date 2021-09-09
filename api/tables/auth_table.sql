CREATE TABLE tokens
(
    Id       	SERIAL PRIMARY KEY,
    NAME			VARCHAR(50),
    Token			VARCHAR(50),
    Scopes 		text[],
    ExpiresAt TIMESTAMP,
    CreatedAt TIMESTAMP,
    Modified  TIMESTAMP
);
