----------------------------------------------------
--          listaway.list table
----------------------------------------------------
CREATE TABLE IF NOT EXISTS listaway.list (
    id SERIAL PRIMARY KEY,
    userid BIGINT NOT NULL,
    name VARCHAR NOT NULL,
    description VARCHAR NULL,
    sharecode VARCHAR
);

CREATE INDEX IF NOT EXISTS list_userid_idx ON listaway.list (userid);
CREATE INDEX IF NOT EXISTS list_sharecode_idx ON listaway.list (sharecode);

----------------------------------------------------
--          listaway.item table
----------------------------------------------------
CREATE TABLE IF NOT EXISTS listaway.item (
    id SERIAL PRIMARY KEY,
    listid BIGINT NOT NULL,
    name VARCHAR NOT NULL,
    url VARCHAR,
    notes VARCHAR,
    priority INT
);

CREATE INDEX IF NOT EXISTS item_listid_idx ON listaway.item (listid);

----------------------------------------------------
--          listaway.user table
----------------------------------------------------
CREATE TABLE IF NOT EXISTS listaway.user (
    id SERIAL,
    email VARCHAR PRIMARY KEY,
    name VARCHAR,
    passwordhash VARCHAR NOT NULL,
    admin BOOLEAN NOT NULL,
    instanceAdmin BOOLEAN NOT NULL
);

-- Migration from 1.6.0 to 1.7.0
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_schema = 'listaway'
        AND table_name = 'user'
        AND column_name = 'instanceadmin'
    ) THEN
        ALTER TABLE listaway.user
        ADD COLUMN InstanceAdmin BOOLEAN DEFAULT false NOT NULL;
    END IF;
    IF NOT EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_schema = 'listaway'
        AND table_name = 'user'
        AND column_name = 'groupid'
    ) THEN
        ALTER TABLE listaway.user
        ADD COLUMN groupid INTEGER DEFAULT 0 NOT NULL; --temporary default to get the column created
        
        ALTER TABLE listaway.user
        ALTER COLUMN groupid DROP DEFAULT;
    END IF;
    
END $$;

CREATE INDEX IF NOT EXISTS user_groupid_idx ON listaway.user (groupid)