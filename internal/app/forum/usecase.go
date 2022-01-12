package forum

import "github.com/perlinleo/technopark-mail.ru-forum-database/internal/model"

type Usecase interface {
	CreateForum(*model.Forum) (*model.Forum, int, error)
	Find(slug string) (*model.Forum, error)
	CreateThread(string, *model.NewThread) (*model.Thread, int, error)
	GetUsersByForum(forumSlug string, params map[string][]string) ([]model.User, int, error)
	GetThreadsByForum(forumSlug string, params map[string][]string) ([]model.Thread, int, error)
}