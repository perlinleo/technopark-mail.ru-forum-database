INSERT INTO forums (slug, title, usernick)
VALUES ($1, $2, $3) RETURNING slug, title, usernick, posts, threads