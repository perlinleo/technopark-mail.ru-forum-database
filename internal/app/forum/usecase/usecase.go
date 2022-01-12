package usecase

import (
	"github.com/perlinleo/technopark-mail.ru-forum-database/internal/app/forum"
	"github.com/perlinleo/technopark-mail.ru-forum-database/internal/model"
)

type ForumUsecase struct {
	repository forum.Repository
}

func NewUserUsecase(r forum.Repository) forum.Usecase {
	return &ForumUsecase{
		repository: r,
	}
}

func (f ForumUsecase) CreateForum(*model.Forum) (*model.Forum, int, error){

}
func (f ForumUsecase) Find(slug string) (*model.Forum, error){

}
func (f ForumUsecase) CreateThread(string, *model.NewThread) (*model.Thread, int, error) {

}
func (f ForumUsecase) GetUsersByForum(forumSlug string, params map[string][]string) ([]model.User, int, error) {

}
func (f ForumUsecase) GetThreadsByForum(forumSlug string, params map[string][]string) ([]model.Thread, int, error) {
	
}
