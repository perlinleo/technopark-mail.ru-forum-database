package repository

import (
	"github.com/jackc/pgx"
	cache "github.com/patrickmn/go-cache"
	"github.com/perlinleo/technopark-mail.ru-forum-database/internal/app/forum"
	"github.com/perlinleo/technopark-mail.ru-forum-database/internal/model"
)

type ForumPSQL struct {
	Conn *pgx.ConnPool
	Cache *cache.Cache
}

func NewForumPSQLRepository(ConnectionPool *pgx.ConnPool, Cache *cache.Cache) forum.Repository {
	return &ForumPSQL{
			ConnectionPool, 
			Cache}
}


func (r ForumPSQL) Create(forum *model.Forum) error{

}
func (r ForumPSQL) Find(slug string) (*model.Forum, error) {

}
func (r ForumPSQL) FindForumUsers(forumObj *model.Forum, params map[string][]string) ([]model.User, error) {

}
func (r ForumPSQL) FindForumThreads(forumSlug string, params map[string][]string) ([]model.Thread, error) {

}