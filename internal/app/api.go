package app

import (
	"fmt"

	"github.com/perlinleo/technopark-mail.ru-forum-database/internal/app/user/delivery"
	"github.com/perlinleo/technopark-mail.ru-forum-database/internal/app/user/repository"
	"github.com/perlinleo/technopark-mail.ru-forum-database/internal/app/user/usecase"
	"github.com/perlinleo/technopark-mail.ru-forum-database/internal/middleware"
	"github.com/valyala/fasthttp"
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


	userRepository := repository.NewUserPSQLRepository(ConnPool,nil)

	userUsecase := usecase.NewUserUsecase(userRepository)

	delivery.NewUserHandler(server.Router, userUsecase)

	fmt.Printf("STARTING SERVICE ON PORT %s\n", config.App.Port)
	err = fasthttp.ListenAndServe(config.App.Port,middleware.ReponseMiddlwareAndLogger(server.Router.Handler))
	if err != nil {
		return err;
	}

	return nil;
}