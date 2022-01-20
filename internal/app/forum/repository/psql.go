package forum_psql

import (
	"fmt"

	"github.com/jackc/pgx"
	cache "github.com/patrickmn/go-cache"
	"github.com/perlinleo/technopark-mail.ru-forum-database/internal/app/forum"
	"github.com/perlinleo/technopark-mail.ru-forum-database/internal/model"
)

type ForumPSQL struct {
	Conn  *pgx.ConnPool
	Cache *cache.Cache
}

func NewForumPSQLRepository(ConnectionPool *pgx.ConnPool, Cache *cache.Cache) forum.Repository {
	return &ForumPSQL{
		ConnectionPool,
		Cache}
}

func (r ForumPSQL) FindForumUsers(forumObj *model.Forum, descValue bool, limitValue string, sinceValue string) ([]model.User, error) {
	limit := "100"
	if limitValue != "" {
		limit = limitValue
	}
	sinceConditionSign := ">"
	desc := ""
	if descValue {
		desc = "desc"
		sinceConditionSign = "<"
	}
	since := ""
	if sinceValue != "" {
		since = sinceValue
	}

	users := []model.User{}

	query := "SELECT nickname, email, fullname, about FROM users " +
		"WHERE id IN (SELECT user_id FROM forum_users WHERE forum_id = $1) "
	if since != "" {
		query += fmt.Sprintf(" AND nickname %s '%s' ", sinceConditionSign, since)
	}
	query += fmt.Sprintf(" ORDER BY nickname %s LIMIT %s ", desc, limit)

	rows, err := r.Conn.Query(query, forumObj.ID)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		u := model.User{}
		err := rows.Scan(&u.Nickname, &u.Email, &u.Fullname, &u.About)

		if err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	return users, nil
}

func (r ForumPSQL) Create(forum *model.Forum) error {
	return r.Conn.QueryRow(
		"INSERT INTO forums (slug, title, usernick) "+
			"VALUES ($1, $2, $3) RETURNING slug, title, usernick, posts, threads",
		forum.Slug,
		forum.Title,
		forum.User,
	).Scan(
		&forum.Slug,
		&forum.Title,
		&forum.User,
		&forum.Posts,
		&forum.Threads,
	)
}
func (r ForumPSQL) Find(slug string) (*model.Forum, error) {
	forum := &model.Forum{}

	if err := r.Conn.QueryRow(
		"SELECT slug, title, usernick, posts, threads, id FROM forums WHERE slug = $1",
		slug,
	).Scan(
		&forum.Slug,
		&forum.Title,
		&forum.User,
		&forum.Posts,
		&forum.Threads,
		&forum.ID,
	); err != nil {
		return nil, err
	}

	return forum, nil
}

func (r ForumPSQL) FindForumThreads(forumSlug string, limitValue string, descValue bool, sinceValue string) ([]model.Thread, error) {
	limit := "100"
	if limitValue != "" {
		limit = limitValue
	}
	desc := ""
	conditionSign := ">="
	if descValue {
		desc = "desc"
		conditionSign = "<="
	}
	since := ""
	if sinceValue != "" {
		since = sinceValue
	}

	threads := []model.Thread{}

	query := "SELECT id, forum, author, slug, created, title, message, votes FROM threads WHERE forum = $1 "

	if since != "" {
		query += fmt.Sprintf(" AND created %s '%s' ", conditionSign, since)
	}

	query += fmt.Sprintf(" ORDER BY created %s LIMIT %s", desc, limit)

	rows, err := r.Conn.Query(query, forumSlug)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		t := model.Thread{}
		err := rows.Scan(&t.ID, &t.Forum, &t.Author, &t.Slug, &t.Created, &t.Title, &t.Message, &t.Votes)

		if err != nil {
			return nil, err
		}

		threads = append(threads, t)
	}

	return threads, nil
}
