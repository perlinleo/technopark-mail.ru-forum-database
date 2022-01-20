package user_usecase

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

func (u UserUsecase) CreateUser(user *model.User) ([]model.User, error) {
	err := u.repository.Create(user)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (u UserUsecase) DuplicateUser(user *model.User) ([]model.User, error) {
	FoundedDuplicates, err := u.repository.Find(user.Nickname, user.Email)
	if err != nil {
		return nil, err
	}

	return FoundedDuplicates, nil
}

func (u UserUsecase) Find(nickname string) (*model.User, error) {
	return u.repository.FindByNickname(nickname)
}

func (u UserUsecase) Update(user *model.User) (*model.User, error, int) {
	

	userObj, err := u.repository.FindByNickname(user.Nickname)

	if err != nil || userObj == nil {
		return nil, err, 404
	}

	if user.Email == "" {
		user.Email = userObj.Email
	}
	if user.About == "" {
		user.About = userObj.About
	}
	if user.Fullname == "" {
		user.Fullname = userObj.Fullname
	}

	userObj, err = u.repository.Update(user)

	if err != nil {
		return nil, err, 409
	}
	return userObj, err, 200
}
