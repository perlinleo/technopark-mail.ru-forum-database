INSERT INTO threads (title, author, forum, message, slug, created)
VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, title, author, forum, message, votes, slug, created;






-- prefer


-- created
INSERT INTO threads (title, message, forum, author, slug, created) " +
			"VALUES ($1, $2, $3, $4, $5, $6) RETURNING slug, title, message, forum, author, created, votes, id


INSERT INTO threads (title, message, forum, author, slug) 
VALUES ($1, $2, $3, $4, $5) RETURNING slug, title, message, forum, author, created, votes, id