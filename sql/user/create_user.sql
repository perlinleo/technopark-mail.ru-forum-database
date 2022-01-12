INSERT INTO users (nickname, email, about, fullname) 
VALUES ($1, $2, $3, $4) RETURNING nickname