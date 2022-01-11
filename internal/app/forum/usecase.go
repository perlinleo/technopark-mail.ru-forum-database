package forum

import "../../model"

type Usecase interface {
	NewForum(*model.Forum) (*model.Forum, int, error)
	FindForum(*model.Forum) (*model.Forum, int, error)
	UsersByForum(forumString string, params map[string][]string) ([]model.User, int, error)
	ThreadsByForum(forumString string, params map[string][]string) ([]model.Thread, int, error)

}