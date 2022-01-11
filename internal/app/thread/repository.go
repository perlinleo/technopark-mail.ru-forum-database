package forum

import "../../model"

type Repository interface {
	Create(forum *model.Forum) error
	Find(slug string) (*model.Forum, error)
	FindForumUsers(forumObj *model.Forum, params map[string][]string) ([]model.User, error)
	FindForumThreads(forumSlug string, params map[string][]string) ([]model.Thread, error)
}
