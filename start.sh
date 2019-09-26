#!/bin/bash
set -e

# if these environment variables are not set define them.
# may be this script is ran outside a container.
if test -z "$POSTGRES_DB" 
then
    export POSTGRES_DB=blogdb
    psql -c "CREATE DATABASE blogdb"
fi

if test -z "$POSTGRES_TEST_DB" 
then
    export POSTGRES_TEST_DB=blogdb_test
fi

psql -c "CREATE DATABASE $POSTGRES_TEST_DB"

psql --dbname "$POSTGRES_TEST_DB" <<-EOSQL
    CREATE SEQUENCE public.articles_id_seq;

    CREATE TABLE public.articles (
        id INTEGER DEFAULT NEXTVAL('articles_id_seq'::regclass) PRIMARY KEY,
        title CHARACTER VARYING(255),
        "content" TEXT,
        author CHARACTER VARYING(255)
    );
EOSQL

psql --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE SEQUENCE public.articles_id_seq;

    CREATE TABLE public.articles (
        id INTEGER DEFAULT NEXTVAL('articles_id_seq'::regclass) PRIMARY KEY,
        title CHARACTER VARYING(255),
        "content" TEXT,
        author CHARACTER VARYING(255)
    );

    INSERT INTO articles(title, "content", author) VALUES('Post 1', 'Some interesting content goes here.', 'Farid');
    INSERT INTO articles(title, "content", author) VALUES('Post 2', 'Again some interesting content goes here.', 'Farid');
    INSERT INTO articles(title, "content", author) VALUES('Post 3', 'Yet more interesting content goes here.', 'Farid');
EOSQL