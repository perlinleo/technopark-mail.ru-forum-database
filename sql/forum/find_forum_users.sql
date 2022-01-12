SELECT nickname, email, fullname, about FROM users
WHERE id IN (SELECT user_id FROM forum_users WHERE forum_id = $1)
ORDER BY nickname DESC LIMIT 100;


SELECT nickname, email, fullname, about FROM users
WHERE id IN (SELECT user_id FROM forum_users WHERE forum_id = $1)
AND nickname < 'nickname' 
ORDER BY nickname DESC LIMIT 100;