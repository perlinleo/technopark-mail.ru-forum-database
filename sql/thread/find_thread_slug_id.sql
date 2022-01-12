SELECT slug, title, message, forum, author, created, votes, id FROM threads 
WHERE id=$1 OR (slug=$2 AND slug <> '');