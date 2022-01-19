package app

import (
	"fmt"
	"time"

	"github.com/patrickmn/go-cache"
	user_psql "github.com/perlinleo/technopark-mail.ru-forum-database/internal/app/user/repository"

	forum_psql "github.com/perlinleo/technopark-mail.ru-forum-database/internal/app/forum/repository"
	thread_psql "github.com/perlinleo/technopark-mail.ru-forum-database/internal/app/thread/repository"

	forum_usecase "github.com/perlinleo/technopark-mail.ru-forum-database/internal/app/forum/usecase"
	thread_usecase "github.com/perlinleo/technopark-mail.ru-forum-database/internal/app/thread/usecase"
	user_usecase "github.com/perlinleo/technopark-mail.ru-forum-database/internal/app/user/usecase"
	"github.com/perlinleo/technopark-mail.ru-forum-database/internal/middleware"
	"github.com/valyala/fasthttp"

	forum_http "github.com/perlinleo/technopark-mail.ru-forum-database/internal/app/forum/delivery"
	thread_http "github.com/perlinleo/technopark-mail.ru-forum-database/internal/app/thread/delivery"
	user_http "github.com/perlinleo/technopark-mail.ru-forum-database/internal/app/user/delivery"
)



func Start() error {
	config := NewConfig()

	server, err := NewServer(config)
	if err !=nil {
		return err
	}
	ConnPool , err := NewDataBase(config.App.DatabaseURL)
	if err !=nil {
		return err
	}

	userCache := cache.New(time.Minute,time.Minute)
	userRepository := user_psql.NewUserPSQLRepository(ConnPool,userCache)
	forumCache := cache.New(time.Minute,time.Minute)
	forumRepository := forum_psql.NewForumPSQLRepository(ConnPool,forumCache)
	threadCache := cache.New(time.Minute,time.Minute)
	threadRepository := thread_psql.NewThreadPSQLRepository(ConnPool,threadCache)
	forumUsecaseCache := cache.New(time.Minute,time.Minute)



	threadUsecase := thread_usecase.NewThreadUsecase(threadRepository, userRepository)
	userUsecase := user_usecase.NewUserUsecase(userRepository)
	forumUsecase := forum_usecase.NewForumUsecase(forumRepository,threadRepository,userRepository,forumUsecaseCache)

	user_http.NewUserHandler(server.Router, userUsecase)
	forum_http.NewForumHandler(server.Router, forumUsecase)
	thread_http.NewThreadHandler(server.Router, threadUsecase)

	fmt.Printf("STARTING SERVICE ON PORT %s\n", config.App.Port)
	err = fasthttp.ListenAndServe(config.App.Port,middleware.ReponseMiddlwareAndLogger(server.Router.Handler))
	if err != nil {
		return err;
	}

	return nil;
}