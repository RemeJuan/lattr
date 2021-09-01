CREATE TABLE tweets
(
    Id       	SERIAL PRIMARY KEY,
    NAME			VARCHAR(300),
    Token			VARCHAR(300),
    Scopes 		[]string
    ExpiresAt TIMESTAMP,
    CreatedAt TIMESTAMP,
    Modified  TIMESTAMP
);
