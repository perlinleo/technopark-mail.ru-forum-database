package forum

import "github.com/perlinleo/technopark-mail.ru-forum-database/internal/model"

type Repository interface {
	Create(forum *model.Forum) error
	Find(slug string) (*model.Forum, error)
	FindForumUsers(forumObj *model.Forum, descValue bool, limitValue string,sinceValue string) ([]model.User, error)
	FindForumThreads(forumSlug string, limitValue string, descValue bool, sinceValue string) ([]model.Thread, error)
}