package repository

import (
	"fmt"

	"github.com/jackc/pgx"
	cache "github.com/patrickmn/go-cache"
	"github.com/perlinleo/technopark-mail.ru-forum-database/internal/app/user"
	"github.com/perlinleo/technopark-mail.ru-forum-database/internal/model"
)

type UserPSQL struct {
	Conn *pgx.ConnPool
	Cache *cache.Cache
}

func NewUserPSQLRepository(ConnectionPool *pgx.ConnPool, Cache *cache.Cache) user.Repository {
	return &UserPSQL{
			ConnectionPool, 
			Cache}
}

func (userRep UserPSQL) Create(user *model.User) error {
	query := "INSERT INTO users (nickname, email, about, fullname) "+
			"VALUES ($1, $2, $3, $4) RETURNING nickname"
	
	return userRep.Conn.QueryRow(
		query, user.Nickname, user.Email, user.About, user.Fullname,).Scan(&user.Nickname)
}


func (userRep UserPSQL) FindByNickname(nickname string) (*model.User, error) {
	userObj := &model.User{}
	fmt.Println("DSIDJHSIUJDISJD")

	if err := userRep.Conn.QueryRow(
			"SELECT nickname, about, email, fullname FROM users WHERE nickname = $1",
			nickname,
		).Scan(
			&userObj.Nickname,
			&userObj.About,
			&userObj.Email,
			&userObj.Fullname,
		); err != nil {
			fmt.Println(err)
			return nil, err
		}

	fmt.Println(userObj)
	return userObj, nil
}


func (u UserPSQL) Update(user *model.User) (*model.User, error) {
	if err := u.Conn.QueryRow(
		"UPDATE users SET about = $1, email = $2, fullname = $3 WHERE nickname = $4 RETURNING nickname, about, email, fullname",
		user.About,
		user.Email,
		user.Fullname,
		user.Nickname,
	).Scan(
		&user.Nickname,
		&user.About,
		&user.Email,
		&user.Fullname,
	); err != nil {
		return nil, err
	}
	return user, nil
}

func (u UserPSQL) Find(nickname string, email string) ([]model.User, error){
	return nil,nil
}
