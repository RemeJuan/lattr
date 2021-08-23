CREATE TABLE tweets
(
    Id       	SERIAL PRIMARY KEY,
    UserId  	VARCHAR(300),
    Message   VARCHAR(300),
    PostTime  TIMESTAMP,
    Status    VARCHAR(10) ,
    CreatedAt TIMESTAMP,
    Modified  TIMESTAMP
);
