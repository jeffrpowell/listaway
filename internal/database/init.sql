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
    instanceAdmin BOOLEAN NOT NULL,
    oidc_provider VARCHAR NULL,
    oidc_subject VARCHAR NULL,
    oidc_email VARCHAR NULL
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

-- Migration from 1.14.0 to 1.15.0 to add OIDC support to user table
DO $$
BEGIN
    -- Add OIDC provider column
    IF NOT EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_schema = 'listaway'
        AND table_name = 'user'
        AND column_name = 'oidc_provider'
    ) THEN
        ALTER TABLE listaway.user
        ADD COLUMN oidc_provider VARCHAR NULL;
    END IF;
    
    -- Add OIDC subject column
    IF NOT EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_schema = 'listaway'
        AND table_name = 'user'
        AND column_name = 'oidc_subject'
    ) THEN
        ALTER TABLE listaway.user
        ADD COLUMN oidc_subject VARCHAR NULL;
    END IF;
    
    -- Add OIDC email column (may differ from primary email)
    IF NOT EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_schema = 'listaway'
        AND table_name = 'user'
        AND column_name = 'oidc_email'
    ) THEN
        ALTER TABLE listaway.user
        ADD COLUMN oidc_email VARCHAR NULL;
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS user_groupid_idx ON listaway.user (groupid);
CREATE INDEX IF NOT EXISTS user_oidc_provider_subject_idx ON listaway.user (oidc_provider, oidc_subject);
CREATE INDEX IF NOT EXISTS user_oidc_email_idx ON listaway.user (oidc_email);

----------------------------------------------------
--          listaway.reset_tokens table
----------------------------------------------------
CREATE TABLE IF NOT EXISTS listaway.reset_tokens (
    token VARCHAR PRIMARY KEY,
    email VARCHAR NOT NULL,
    created_at TIMESTAMP NOT NULL,
    expires_at TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS reset_tokens_email_idx ON listaway.reset_tokens (email);

----------------------------------------------------
--          listaway.collection table
----------------------------------------------------
CREATE TABLE IF NOT EXISTS listaway.collection (
    id SERIAL PRIMARY KEY,
    userid BIGINT NOT NULL,
    name VARCHAR NOT NULL,
    description VARCHAR NULL,
    sharecode VARCHAR
);

CREATE INDEX IF NOT EXISTS collection_userid_idx ON listaway.collection (userid);
CREATE INDEX IF NOT EXISTS collection_sharecode_idx ON listaway.collection (sharecode);

----------------------------------------------------
--          listaway.collection_list table
----------------------------------------------------
CREATE TABLE IF NOT EXISTS listaway.collection_list (
    collectionid BIGINT NOT NULL,
    listid BIGINT NOT NULL,
    PRIMARY KEY (collectionid, listid)
);

CREATE INDEX IF NOT EXISTS collection_list_collectionid_idx ON listaway.collection_list (collectionid);
CREATE INDEX IF NOT EXISTS collection_list_listid_idx ON listaway.collection_list (listid);
