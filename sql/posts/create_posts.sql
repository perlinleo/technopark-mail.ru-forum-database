INSERT INTO posts(id, parent, thread, forum, author, created, message, path) 
VALUES (nextval('posts_id_seq'::regclass), ?, ?, ?, ?, ?, ?,
				ARRAY[currval(pg_get_serial_sequence('posts', 'id'))::bigint])
                 RETURNING  id, parent, thread, forum, author, created, message, isedited 

SELECT thread FROM posts WHERE id = $1

-- for
SELECT nickname FROM users WHERE nickname = $1 

-- next
(nextval('posts_id_seq'::regclass), ?, ?, ?, ?, ?, ?, " +
				"(SELECT path FROM posts WHERE id = ? AND thread = ?) || " +
				"currval(pg_get_serial_sequence('posts', 'id'))::bigint)