SET SYNCHRONOUS_COMMIT = 'off';
create extension if not exists citext;


DROP INDEX IF EXISTS idx_users_email_uindex;
DROP INDEX IF EXISTS idx_users_nickname_uindex;
DROP INDEX IF EXISTS idx_forums_slug_uindex;
DROP INDEX IF EXISTS idx_forums_userNick_unique;;
DROP INDEX IF EXISTS idx_threads_slug;
DROP INDEX IF EXISTS idx_threads_forum;
DROP INDEX IF EXISTS idx_posts_forum;
DROP INDEX IF EXISTS idx_posts_parent;
DROP INDEX IF EXISTS idx_posts_path;
DROP INDEX IF EXISTS idx_posts_thread;
DROP INDEX IF EXISTS idx_posts_thread_id;

DROP TRIGGER IF EXISTS on_vote_insert ON votes;
DROP TRIGGER IF EXISTS on_vote_update ON votes;

DROP FUNCTION IF EXISTS fn_update_thread_votes_ins();
DROP FUNCTION IF EXISTS fn_update_thread_votes_upd();

DROP TABLE IF EXISTS forum_users;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS forums;
DROP TABLE IF EXISTS threads CASCADE;
DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS votes;


CREATE TABLE IF NOT EXISTS users
(
    id       serial not null primary key,
    nickname varchar(20) not null unique,
    about    varchar(255),
    email    varchar(20)  not null,
    fullname varchar(50)   not null
);

-- CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email_uindex
--     ON users (LOWER(email));
-- CREATE UNIQUE INDEX IF NOT EXISTS idx_users_nickname_uindex
--     ON users (LOWER(nickname));
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email_uindex2
    ON users (email);
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_nickname_uindex2
    ON users (nickname);
CREATE INDEX IF NOT EXISTS idx_users_pok
    ON users (nickname, email, fullname, about);

CREATE TABLE IF NOT EXISTS forums
(
    id       serial not null primary key,
    slug     citext   not null,
    userNick citext   not null,
    title    varchar,
    posts    int default 0,
    threads  int default 0
);

-- CREATE UNIQUE INDEX IF NOT EXISTS idx_forums_slug_uindex
--     ON forums (LOWER(slug));
-- CREATE UNIQUE INDEX IF NOT EXISTS idx_forums_userNick_unique
--     ON forums (LOWER(userNick));
CREATE UNIQUE INDEX IF NOT EXISTS idx_forums_slug_uindex2
    ON forums (slug);
CREATE UNIQUE INDEX IF NOT EXISTS idx_forums_userNick_unique2
    ON forums (userNick);

CREATE TABLE IF NOT EXISTS threads
(
    id      serial not null primary key,
    slug    citext,
    title   varchar,
    message varchar,
    votes   int         default 0,
    author  varchar,
    forum   citext,
    created timestamptz DEFAULT now()
);

-- CREATE INDEX IF NOT EXISTS idx_threads_slug
--     ON threads (LOWER(slug));
-- CREATE INDEX IF NOT EXISTS idx_threads_forum
--     ON threads (LOWER(forum));
CREATE INDEX IF NOT EXISTS idx_threads_slug2
    ON threads (slug);
CREATE INDEX IF NOT EXISTS idx_threads_forum2
    ON threads (forum);
CREATE INDEX IF NOT EXISTS idx_threads_pok
    ON threads (id, forum, author, slug, created, title, message, votes);
CREATE INDEX IF NOT EXISTS idx_threads_created
    ON threads (created);
CREATE INDEX IF NOT EXISTS idx_threads_created2
    ON threads (created, forum);

CREATE TABLE IF NOT EXISTS posts
(
    id       serial not null primary key,
    parent   int             DEFAULT NULL,
    path     int[]  NOT NULL DEFAULT '{0}',
    thread   int REFERENCES threads(id) NOT NULL,
    forum    citext,
    author   citext,
    created  timestamptz        DEFAULT now(),
    isEdited bool               DEFAULT FALSE,
    message  text
);
CREATE INDEX IF NOT EXISTS idx_posts_path ON posts USING GIN (path);
CREATE INDEX IF NOT EXISTS idx_posts_thread ON posts (thread);
CREATE INDEX IF NOT EXISTS idx_posts_forum ON posts (forum);
CREATE INDEX IF NOT EXISTS idx_posts_parent ON posts (parent);
CREATE INDEX IF NOT EXISTS idx_posts_thread_id ON posts (thread, id);
CREATE INDEX IF NOT EXISTS idx_posts_pok
    ON posts (id, parent, thread, forum, author, created, message, isedited, path);
CREATE INDEX IF NOT EXISTS idx_posts_created
    ON posts (created);

CREATE TABLE IF NOT EXISTS votes
(
    nickname citext  REFERENCES users(nickname) NOT NULL,
    thread   int      REFERENCES threads(id) NOT NULL,
    voice    smallint NOT NULL CHECK (voice = 1 OR voice = -1),
    PRIMARY KEY (nickname, thread)
);
-- CREATE UNIQUE INDEX IF NOT EXISTS idx_votes_nickname_thread_unique
--     ON votes (LOWER(nickname), thread);
CREATE UNIQUE INDEX IF NOT EXISTS idx_votes_nickname_thread_unique2
    ON votes (nickname, thread);

CREATE FUNCTION fn_update_thread_votes_ins()
    RETURNS TRIGGER AS '
    BEGIN
        UPDATE threads
        SET
            votes = votes + NEW.voice
        WHERE id = NEW.thread;
        RETURN NULL;
    END;
' LANGUAGE plpgsql;


CREATE TRIGGER on_vote_insert
    AFTER INSERT ON votes
    FOR EACH ROW EXECUTE PROCEDURE fn_update_thread_votes_ins();

CREATE FUNCTION fn_update_thread_votes_upd()
    RETURNS TRIGGER AS '
    BEGIN
        IF OLD.voice = NEW.voice
        THEN
            RETURN NULL;
        END IF;
        UPDATE threads
        SET
            votes = votes + CASE WHEN NEW.voice = -1
                                     THEN -2
                                 ELSE 2 END
        WHERE id = NEW.thread;
        RETURN NULL;
    END;
' LANGUAGE plpgsql;

CREATE TRIGGER on_vote_update
    AFTER UPDATE ON votes
    FOR EACH ROW EXECUTE PROCEDURE fn_update_thread_votes_upd();


CREATE TABLE IF NOT EXISTS forum_users (
    user_id INT REFERENCES users(id),
    forum_id INT REFERENCES forums(id)
);

CREATE INDEX idx_forum_users_user_id
    ON forum_users(user_id);

CREATE INDEX idx_forum_users_forum_id
    ON forum_users(forum_id);

CREATE INDEX idx_forum_users_user_id_forum_id
    ON forum_users (user_id, forum_id);

CREATE OR REPLACE FUNCTION forum_users_update()
    RETURNS TRIGGER AS '
    BEGIN
--         INSERT INTO forum_users (user_id, forum_id) VALUES ((SELECT id FROM users WHERE LOWER(NEW.author) = LOWER(nickname)),
--                                                               (SELECT id FROM forums WHERE LOWER(NEW.forum) = LOWER(slug)));
        INSERT INTO forum_users (user_id, forum_id) VALUES ((SELECT id FROM users WHERE NEW.author = nickname),
                                                            (SELECT id FROM forums WHERE NEW.forum = slug));
        RETURN NULL;
    END;
' LANGUAGE plpgsql;

CREATE TRIGGER on_post_insert
    AFTER INSERT ON posts
    FOR EACH ROW EXECUTE PROCEDURE forum_users_update();

CREATE TRIGGER on_thread_insert
    AFTER INSERT ON threads
    FOR EACH ROW EXECUTE PROCEDURE forum_users_update();


CLUSTER users USING users_nickname_key;
CLUSTER threads USING idx_threads_created2;
CLUSTER forums USING idx_forums_slug_uindex2;
CLUSTER posts USING idx_posts_thread_id;
