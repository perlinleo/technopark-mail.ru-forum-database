SELECT
(SELECT COUNT(*) FROM forums) AS forum, 
(SELECT COUNT(*) FROM posts) AS post, 
(SELECT COUNT(*) FROM threads) AS thread, 
(SELECT COUNT(*) FROM users) AS user;