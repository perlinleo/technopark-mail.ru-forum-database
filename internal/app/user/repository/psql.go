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

func (userRep *UserPSQL) Create(user *model.User) error {
	query := "INSERT INTO users (nickname, email, about, fullname) "+
			"VALUES ($1, $2, $3, $4) RETURNING nickname"
	fmt.Printf("%v\n\n",user.Nickname)
	_, err := userRep.Conn.Exec(
		query, user.Nickname, user.Email, user.About, user.Fullname,)
	
	return err
}


func (userRep *UserPSQL) FindByNickname(nickname string) (*model.User, error) {
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
			fmt.Println(err)
			return nil, err
		}

	return userObj, nil
}


func (u *UserPSQL) Update(user *model.User) (*model.User, error) {
	_, err := u.Conn.Exec(
		"UPDATE users SET about = $1, email = $2, fullname = $3 WHERE nickname = $4 RETURNING nickname, about, email, fullname",
		user.About,
		user.Email,
		user.Fullname,
		user.Nickname,
	)
	if err != nil {
		return nil, err
	}
	
	return user, nil
}

func (u *UserPSQL) Find(nickname string, email string) ([]model.User, error){
	var users []model.User
	rows,err := u.Conn.Query(`SELECT nickname, about,email, fullname FROM users 
		WHERE nickname = $1 OR email = $2`, nickname,email)
	if err != nil {
		fmt.Println(err)
	}
	for rows.Next(){
		obj := model.User{}
		rows.Scan(&obj.Nickname, &obj.About, &obj.Email, &obj.Fullname)
		users = append(users, obj)
	}
	return users,nil
}
