CREATE TABLE IF NOT EXISTS listaway.list (
    Id SERIAL PRIMARY KEY,
    UserId BIGINT NOT NULL,
    Name VARCHAR NOT NULL,
    ShareCode VARCHAR
);

CREATE INDEX ON listaway.list (UserId);
CREATE INDEX ON listaway.list (ShareCode);

CREATE TABLE IF NOT EXISTS listaway.item (
    Id SERIAL PRIMARY KEY,
    ListId BIGINT NOT NULL,
    Name VARCHAR NOT NULL,
    URL VARCHAR,
    Notes VARCHAR,
    Priority INT
);

CREATE INDEX ON listaway.item (ListId);

CREATE TABLE IF NOT EXISTS listaway.user (
    Id SERIAL,
    Email VARCHAR PRIMARY KEY,
    Name VARCHAR,
    PasswordHash VARCHAR NOT NULL,
    Admin BOOLEAN NOT NULL
);