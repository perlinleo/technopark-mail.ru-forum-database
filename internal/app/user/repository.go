package user

import "github.com/perlinleo/technopark-mail.ru-forum-database/internal/model"

type Repository interface {
	Create(forum *model.User) error
	FindByNickname(nickname string) (*model.User, error)
	Find(nickname string, email string) ([]model.User, error)
	Update(user *model.User) (*model.User, error)
}
