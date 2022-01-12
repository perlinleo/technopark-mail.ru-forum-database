package user

import "github.com/perlinleo/technopark-mail.ru-forum-database/internal/model"

type Usecase interface {
	CreateUser(user *model.User) ([]model.User, error)
	DuplicateUser(user *model.User) ([]model.User, error)
	Find(nickname string) (*model.User, error)
	Update(user *model.User) (*model.User, error, int)
}
