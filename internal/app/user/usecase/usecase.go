package usecase

import (
	"github.com/perlinleo/technopark-mail.ru-forum-database/internal/app/user"
	"github.com/perlinleo/technopark-mail.ru-forum-database/internal/model"
)

type UserUsecase struct {
	repository user.Repository
}

func NewUserUsecase(r user.Repository) user.Usecase {
	return &UserUsecase{
		repository: r,
	}
}

func (u UserUsecase) CreateUser(user *model.User) (error) {
	err := u.repository.Create(user)
	if err != nil {
		return err
	}
	return nil
}


func (u UserUsecase) Find(nickname string) (*model.User, error) {
	return u.repository.FindByNickname(nickname)
}


func (u UserUsecase) Update(user *model.User) (*model.User, error, int) {
	userObj, err := u.repository.FindByNickname(user.Nickname)

	if err != nil || userObj == nil {
		return nil, err, 404
	}

	userObj, err = u.repository.Update(user)

	if err != nil {
		return nil, err, 409
	}

	return userObj, err, 200
}

