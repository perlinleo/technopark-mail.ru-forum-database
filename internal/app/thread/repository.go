package thread

import "github.com/perlinleo/technopark-mail.ru-forum-database/internal/model"

type Repository interface {
	CreateThread(thread *model.NewThread) (*model.Thread, error)
	FindByIdOrSlug(id int, slug string) (*model.Thread, error)
	CreatePosts(thread *model.Thread, posts []*model.Post) ([]*model.Post, error)
	UpdateThread(id int, slug string, update *model.ThreadUpdate) (*model.Thread, error)
	GetThreadPosts(thread *model.Thread, limit, desc, since, sort string) ([]model.Post, error)
	Vote(thread *model.Thread, vote *model.Vote) (*model.Thread, error)
}

