CREATE TABLE IF NOT EXISTS listaway.list (
    Id SERIAL PRIMARY KEY,
    UserId BIGINT,
    Name NVARCHAR,
    ShareURL NVARCHAR,
    INDEX UserIdIndex (UserId),
    INDEX ShareURLIndex (ShareURL)
);


CREATE TABLE IF NOT EXISTS listaway.item (
    Id SERIAL PRIMARY KEY,
    ListId BIGINT,
    Name NVARCHAR,
    URL NVARCHAR,
    Notes NVARCHAR,
    Priority INT,
    INDEX (ListId)
);

CREATE TABLE IF NOT EXISTS listaway.user (
    Id SERIAL,
    Email NVARCHAR PRIMARY KEY,
    Name NVARCHAR,
    PasswordHash NVARCHAR
);
