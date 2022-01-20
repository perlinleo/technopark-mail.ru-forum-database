package forum_usecase

import (
	"github.com/patrickmn/go-cache"
	"github.com/perlinleo/technopark-mail.ru-forum-database/internal/app/forum"
	"github.com/perlinleo/technopark-mail.ru-forum-database/internal/app/thread"
	"github.com/perlinleo/technopark-mail.ru-forum-database/internal/app/user"
	"github.com/perlinleo/technopark-mail.ru-forum-database/internal/model"
)

type ForumUsecase struct {
	repositoryForum  forum.Repository
	repositoryUser   user.Repository
	repositoryThread thread.Repository
	Cache            *cache.Cache
}

func NewForumUsecase(repositoryForum forum.Repository, repositoryThread thread.Repository, repositoryUser user.Repository, cache *cache.Cache) forum.Usecase {
	return &ForumUsecase{
		repositoryForum:  repositoryForum,
		repositoryThread: repositoryThread,
		repositoryUser:   repositoryUser,
		Cache:            cache,
	}
}

func (f ForumUsecase) CreateForum(data *model.Forum) (*model.Forum, int, error) {
	userObj, err := f.repositoryUser.FindByNickname(data.User)
	if userObj == nil || err != nil {
		return nil, 404, err
	}

	data.User = userObj.Nickname

	if err := f.repositoryForum.Create(data); err != nil {
		forumObj, err := f.repositoryForum.Find(data.Slug)
		if err != nil {
			return nil, 409, err
		}

		return forumObj, 409, err
	}

	return data, 201, nil
}
func (f ForumUsecase) Find(slug string) (*model.Forum, error) {
	var forumObj *model.Forum
	var err error
	forumObj, err = f.repositoryForum.Find(slug)
	if err != nil {
		return nil, err
	}
	f.Cache.Set(slug, forumObj, cache.DefaultExpiration)
	return forumObj, nil
}
func (f ForumUsecase) CreateThread(slug string, newThread *model.NewThread) (*model.Thread, int, error) {
	userObj, err := f.repositoryUser.FindByNickname(newThread.Author)
	if userObj == nil || err != nil {
		return nil, 404, err
	}

	forumObj, err := f.Find(slug)
	if forumObj == nil || err != nil {
		return nil, 404, err
	}

	newThread.Forum = forumObj.Slug

	threadObj, err := f.repositoryThread.FindByIdOrSlug(0, newThread.Slug)
	if threadObj != nil {
		return threadObj, 409, err
	}

	threadObj, err = f.repositoryThread.CreateThread(newThread)
	if err != nil {
		return nil, 409, err
	}

	return threadObj, 201, nil
}
func (f ForumUsecase) GetUsersByForum(slug string, limitValue string, descValue bool, sinceValue string) ([]model.User, int, error) {
	forumObj, err := f.Find(slug)
	if forumObj == nil || err != nil {
		return nil, 404, err
	}

	users, err := f.repositoryForum.FindForumUsers(forumObj, descValue, limitValue, sinceValue)
	if err != nil {
		return nil, 404, err
	}

	return users, 200, nil
}
func (f ForumUsecase) GetThreadsByForum(forumSlug string, limitValue string, descValue bool, sinceValue string) ([]model.Thread, int, error) {
	forumObj, err := f.Find(forumSlug)
	if forumObj == nil || err != nil {
		return nil, 404, err
	}

	threads, err := f.repositoryForum.FindForumThreads(forumSlug, limitValue, descValue, sinceValue)
	if err != nil {
		return nil, 404, err
	}

	return threads, 200, nil
}
