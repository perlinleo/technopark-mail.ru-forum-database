package thread

import "github.com/perlinleo/technopark-mail.ru-forum-database/internal/model"

type Usecase interface {
	CreatePosts(threadSlugOrId string, posts []*model.Post) ([]*model.Post, int, error)
	FindByIdOrSlug(threadSlugOrId string) (*model.Thread, error)
	UpdateThread(threadSlugOrId string, update *model.ThreadUpdate) (*model.Thread, error)
	GetThreadPosts(threadSlugOrId string, limit string, desc bool, since string, sort string) ([]model.Post, error)
	Vote(threadSlugOrId string, vote *model.Vote) (*model.Thread, error)
}