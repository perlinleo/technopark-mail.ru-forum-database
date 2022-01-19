package thread_psql

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx"
	cache "github.com/patrickmn/go-cache"
	"github.com/perlinleo/technopark-mail.ru-forum-database/internal/app/thread"
	"github.com/perlinleo/technopark-mail.ru-forum-database/internal/model"
)

type ThreadPSQL struct {
	Conn *pgx.ConnPool
	Cache *cache.Cache
}


func NewThreadPSQLRepository(ConnectionPool *pgx.ConnPool, Cache *cache.Cache) thread.Repository {
	return &ThreadPSQL{
			ConnectionPool, 
			Cache}
}


func (t ThreadPSQL) Vote(thread *model.Thread, vote *model.Vote) (*model.Thread, error) {
	_ ,err := t.Conn.Exec("INSERT INTO votes(nickname, voice, thread) VALUES ($1, $2, $3) ON CONFLICT ON CONSTRAINT votes_pkey DO UPDATE SET voice = $2",
	 vote.Nickname, vote.Voice, thread.ID)
	if err != nil {
		return nil, err
	}
	var value int32;
	err = t.Conn.QueryRow("SELECT votes FROM threads WHERE id = $1", thread.ID).Scan(	&value )
	thread.Votes = value;
	t.Cache.Delete(thread.Slug)
	return thread , nil
}

func (t ThreadPSQL) GetThreadPosts(thread *model.Thread, limit, desc, since, sort string) ([]model.Post, error) {
	posts := make([]model.Post, 0)

	var query string

	conditionSign := ">"
	if desc == "desc" {
		conditionSign = "<"
	}

	if sort == "flat" {
		query = "SELECT id, parent, thread, forum, author, created, message, isedited FROM posts WHERE thread = $1 "
		if since != "" {
			query += fmt.Sprintf(" AND id %s %s ", conditionSign, since)
		}
		query += fmt.Sprintf(" ORDER BY created %s, id %s LIMIT %s", desc, desc, limit)
	} else if sort == "tree" {
		orderString := fmt.Sprintf(" ORDER BY path[1] %s, path %s ", desc, desc)

		query = "SELECT id, parent, thread, forum, author, created, message, isedited " +
			"FROM posts " +
			"WHERE thread = $1 "
		if since != "" {
			query += fmt.Sprintf(" AND path %s (SELECT path FROM posts WHERE id = %s) ", conditionSign, since)
		}
		query += orderString
		query += fmt.Sprintf("LIMIT %s", limit)

	} else if sort == "parent_tree" {
		query = "SELECT id, parent, thread, forum, author, created, message, isedited " +
			"FROM posts " +
			"WHERE thread = $1 AND path && (SELECT ARRAY (select id from posts WHERE thread = $1 AND parent = 0 "
		if since != "" {
			query += fmt.Sprintf(" AND path %s (SELECT path[1:1] FROM posts WHERE id = %s) ", conditionSign, since)
		}
		query += fmt.Sprintf("ORDER BY path[1] %s, path LIMIT %s)) ", desc, limit)
		query += fmt.Sprintf("ORDER BY path[1] %s, path ", desc)
	}

	rows, err := t.Conn.Query(query, thread.ID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		p := model.Post{}
		err := rows.Scan(&p.ID, &p.Parent, &p.Thread, &p.Forum, &p.Author, &p.Created, &p.Message, &p.IsEdited)
		if err != nil {
			return nil, err
		}

		posts = append(posts, p)
	}

	//t.cache.Set(key, posts, cache2.DefaultExpiration)
	//}

	return posts, nil
}

func (t ThreadPSQL) UpdateThread(id int, slug string, threadUpdate *model.ThreadUpdate) (*model.Thread, error) {
	th := &model.Thread{}

	err := t.Conn.QueryRow(
		"UPDATE threads SET title = $1, message = $2 WHERE id=$3 OR slug=$4 RETURNING slug, title, message, forum, author, created, votes, id",
		threadUpdate.Title,
		threadUpdate.Message,
		id,
		slug,
	).Scan(
		&th.Slug,
		&th.Title,
		&th.Message,
		&th.Forum,
		&th.Author,
		&th.Created,
		&th.Votes,
		&th.ID,
	)

	if err != nil {
		return nil, err
	}

	t.Cache.Set(th.Slug, th, cache.DefaultExpiration)
	//t.cache.Delete("thread_" + fmt.Sprint(th.ID))

	return th, nil
}

func (t ThreadPSQL) CreatePosts(thread *model.Thread, posts []*model.Post) ([]*model.Post, error) {
	tx, err := t.Conn.Begin()
	if err != nil {
		return nil, err
	}

	now := time.Now()

	sqlStr := "INSERT INTO posts(id, parent, thread, forum, author, created, message, path) VALUES "
	vals := []interface{}{}
	for _, post := range posts {
		var author string
		err = t.Conn.QueryRow("SELECT nickname FROM users WHERE nickname = $1",
			post.Author,
		).Scan(&author)
		if err != nil || author == "" {
			_ = tx.Rollback()
			return nil, errors.New("404")
		}

		if post.Parent == 0 {
			sqlStr += "(nextval('posts_id_seq'::regclass), ?, ?, ?, ?, ?, ?, " +
				"ARRAY[currval(pg_get_serial_sequence('posts', 'id'))::bigint]),"
			vals = append(vals, post.Parent, thread.ID, thread.Forum, post.Author, now, post.Message)
		} else {
			var parentThreadId int32
			err = t.Conn.QueryRow("SELECT thread FROM posts WHERE id = $1",
				post.Parent,
			).Scan(&parentThreadId)
			if err != nil {
				_ = tx.Rollback()
				return nil, err
			}
			if parentThreadId != thread.ID {
				_ = tx.Rollback()
				return nil, errors.New("Parent post was created in another thread")
			}

			sqlStr += " (nextval('posts_id_seq'::regclass), ?, ?, ?, ?, ?, ?, " +
				"(SELECT path FROM posts WHERE id = ? AND thread = ?) || " +
				"currval(pg_get_serial_sequence('posts', 'id'))::bigint),"

			vals = append(vals, post.Parent, thread.ID, thread.Forum, post.Author, now, post.Message, post.Parent, thread.ID)
		}

	}
	sqlStr = strings.TrimSuffix(sqlStr, ",")

	sqlStr += " RETURNING  id, parent, thread, forum, author, created, message, isedited "

	sqlStr = ReplaceSQL(sqlStr, "?")
	if len(posts) > 0 {
		// mb kringe
		// stmtButch, err := tx.Prepare("name1", sqlStr)
		if err != nil {
			return nil, err
		}
		rows, err := tx.Query(sqlStr, vals...)
		if err != nil {
			_ = tx.Rollback()
			return nil, err
		}
		i := 0
		for rows.Next() {
			err := rows.Scan(
				&(posts)[i].ID,
				&(posts)[i].Parent,
				&(posts)[i].Thread,
				&(posts)[i].Forum,
				&(posts)[i].Author,
				&(posts)[i].Created,
				&(posts)[i].Message,
				&(posts)[i].IsEdited,
			)
			i += 1

			if err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	f := &model.Forum{}
	err = t.Conn.QueryRow(
		"UPDATE forums SET posts = posts + $1 WHERE slug = $2 RETURNING slug, title, usernick, posts, threads, id",
		len(posts),
		thread.Forum,
	).Scan(
		&f.Slug,
		&f.Title,
		&f.User,
		&f.Posts,
		&f.Threads,
		&f.ID,
	)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	t.Cache.Set(thread.Forum, f, cache.DefaultExpiration)

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (t ThreadPSQL) CreateThread(newThread *model.NewThread) (*model.Thread, error) {
	th := &model.Thread{}
	var row *pgx.Row

	if newThread.Created == "" {
		query := "INSERT INTO threads (title, message, forum, author, slug) " +
			"VALUES ($1, $2, $3, $4, $5) RETURNING slug, title, message, forum, author, created, votes, id"
		row = t.Conn.QueryRow(
			query,
			newThread.Title,
			newThread.Message,
			newThread.Forum,
			newThread.Author,
			newThread.Slug,
		)
	} else {
		 
		query := "INSERT INTO threads (title, message, forum, author, slug, created) " +
			"VALUES ($1, $2, $3, $4, $5, $6) RETURNING slug, title, message, forum, author, created, votes, id"
		row = t.Conn.QueryRow(
			query,
			newThread.Title,
			newThread.Message,
			newThread.Forum,
			newThread.Author,
			newThread.Slug,
			newThread.Created,
		)
	}

	err := row.Scan(
		&th.Slug,
		&th.Title,
		&th.Message,
		&th.Forum,
		&th.Author,
		&th.Created,
		&th.Votes,
		&th.ID,
	)
	if err != nil {
		return nil, err
	}

	f := &model.Forum{}
	err = t.Conn.QueryRow(
		"UPDATE forums SET threads = threads + 1 WHERE slug = $1 RETURNING slug, title, usernick, posts, threads, id",
		th.Forum,
	).Scan(
		&f.Slug,
		&f.Title,
		&f.User,
		&f.Posts,
		&f.Threads,
		&f.ID,
	)
	if err != nil {
		return nil, err
	}

	t.Cache.Set(th.Forum, f, cache.DefaultExpiration)

	return th, nil
}

func (t ThreadPSQL) FindByIdOrSlug(id int, slug string) (*model.Thread, error) {
	th := &model.Thread{}

	if x, found := t.Cache.Get(slug); found {
		th = x.(*model.Thread)
	} else {
		err := t.Conn.QueryRow(
			"SELECT slug, title, message, forum, author, created, votes, id FROM threads WHERE id=$1 OR (slug=$2 AND slug <> '')",
			id,
			slug,
		).Scan(
			&th.Slug,
			&th.Title,
			&th.Message,
			&th.Forum,
			&th.Author,
			&th.Created,
			&th.Votes,
			&th.ID,
		)

		if err != nil {
			return nil, err
		}

		//t.cache.Delete("thread_" + fmt.Sprint(th.ID))
	}

	return th, nil
}



func ReplaceSQL(old, searchPattern string) string {
	tmpCount := strings.Count(old, searchPattern)
	for m := 1; m <= tmpCount; m++ {
		old = strings.Replace(old, searchPattern, "$"+strconv.Itoa(m), 1)
	}
	return old
}

func (t ThreadPSQL) FindPostId(id string, includeUser, includeForum, includeThread bool) (*model.PostFull, error) {
	postObj := &model.PostFull{}
	postObj.Post = &model.Post{}

	id2, _ := strconv.Atoi(id)

	if x, found := t.Cache.Get("post_" + id); found {
		postObj.Post = x.(*model.Post)
	} else {
		if err := t.Conn.QueryRow(
			"SELECT author, created, forum, id, message, thread, isedited, parent FROM posts WHERE id = $1",
			id2,
		).Scan(
			&postObj.Post.Author,
			&postObj.Post.Created,
			&postObj.Post.Forum,
			&postObj.Post.ID,
			&postObj.Post.Message,
			&postObj.Post.Thread,
			&postObj.Post.IsEdited,
			&postObj.Post.Parent,
		); err != nil {
			return nil, err
		}
	}

	if includeUser {
		postObj.Author = &model.User{}

		if x, found := t.Cache.Get(postObj.Post.Author); found {
			postObj.Author = x.(*model.User)
		} else {
			if err := t.Conn.QueryRow(
				"SELECT about, email, fullname, nickname FROM users WHERE nickname = $1",
				postObj.Post.Author,
			).Scan(
				&postObj.Author.About,
				&postObj.Author.Email,
				&postObj.Author.Fullname,
				&postObj.Author.Nickname,
			); err != nil {
				return nil, err
			}
		}
	}

	if includeForum {
		postObj.Forum = &model.Forum{}

		if x, found := t.Cache.Get(postObj.Post.Forum); found {
			postObj.Forum = x.(*model.Forum)
		} else {
			if err := t.Conn.QueryRow(
				"SELECT usernick, title, slug, posts, threads FROM forums WHERE slug = $1",
				postObj.Post.Forum,
			).Scan(
				&postObj.Forum.User,
				&postObj.Forum.Title,
				&postObj.Forum.Slug,
				&postObj.Forum.Posts,
				&postObj.Forum.Threads,
			); err != nil {
				return nil, err
			}
		}
	}

	if includeThread {
		postObj.Thread = &model.Thread{}

		//if x, found := p.cache.Get("thread_"+fmt.Sprint(postObj.Post.Thread)); found {
		//	postObj.Thread = x.(*model.Thread)
		//} else {
		if err := t.Conn.QueryRow(
			"SELECT forum, slug, title, author, message, id, created, votes FROM threads WHERE id = $1",
			postObj.Post.Thread,
		).Scan(
			&postObj.Thread.Forum,
			&postObj.Thread.Slug,
			&postObj.Thread.Title,
			&postObj.Thread.Author,
			&postObj.Thread.Message,
			&postObj.Thread.ID,
			&postObj.Thread.Created,
			&postObj.Thread.Votes,
		); err != nil {
			return nil, err
		}
		//}
	}

	return postObj, nil
}


func (t ThreadPSQL) UpdatePost(id string, message string) (*model.Post, error) {
	postObj := &model.Post{}

	id2, _ := strconv.Atoi(id)

	if err := t.Conn.QueryRow(
		"UPDATE posts SET message = $2, isEdited = TRUE WHERE id = $1 RETURNING author, created, forum, id, message, thread, isEdited",
		id2,
		message,
	).Scan(
		&postObj.Author,
		&postObj.Created,
		&postObj.Forum,
		&postObj.ID,
		&postObj.Message,
		&postObj.Thread,
		&postObj.IsEdited,
	); err != nil {
		return nil, err
	}

	t.Cache.Set("post_"+id, postObj, cache.DefaultExpiration)

	return postObj, nil
}


func (t ThreadPSQL) ClearAll() error {
	if _, err := t.Conn.Exec("TRUNCATE votes, users, posts, threads, forums RESTART IDENTITY CASCADE"); err != nil {
		return err
	}

	return nil
}

func (t ThreadPSQL) GetStatus() (*model.Status, error) {
	status := &model.Status{}

	//mb cringe
	if err := t.Conn.QueryRow("SELECT "+
		"(SELECT count(*) from forums) AS forum, "+
		"(SELECT count(*) from posts) AS post, "+
		"(SELECT count(*) from threads) AS thread, "+
		"(SELECT count(*) from users) AS user",
	).Scan(&status.Forum, &status.Post, &status.Thread, &status.User); err != nil {
		return nil, err
	}

	return status, nil
}
