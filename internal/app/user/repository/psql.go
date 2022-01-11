package repository

import (
	"../../../model"
	"github.com/jackc/pgx"
	cache "github.com/patrickmn/go-cache"
)

type UserRepository struct {
	Conn *pgx.ConnPool
	Cache *cache.Cache
}


func (userRep UserRepository) Create(user *model.User) error {
	query := "INSERT INTO users (nickname, email, about, fullname) "+
			"VALUES ($1, $2, $3, $4) RETURNING nickname"
	
	return userRep.Conn.QueryRow(
		query, user.Nickname, user.Email, user.About, user.Fullname,).Scan(&user.Nickname)
}


func (userRep UserRepository) FindByNickname(nickname string) (*model.User, error) {
	userObj := &model.User{}

	if err := userRep.Conn.QueryRow(
			"SELECT nickname, about, email, fullname FROM users WHERE nickname = $1",
			nickname,
		).Scan(
			&userObj.Nickname,
			&userObj.About,
			&userObj.Email,
			&userObj.Fullname,
		); err != nil {
			return nil, err
		}
	return userObj, nil
}