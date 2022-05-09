package app

import (
	"fmt"
	"time"

	router "github.com/fasthttp/router"
	"github.com/patrickmn/go-cache"
	forum_psql "github.com/perlinleo/technopark-mail.ru-forum-database/internal/app/forum/repository"
	thread_psql "github.com/perlinleo/technopark-mail.ru-forum-database/internal/app/thread/repository"
	user_psql "github.com/perlinleo/technopark-mail.ru-forum-database/internal/app/user/repository"
	responses "github.com/perlinleo/technopark-mail.ru-forum-database/internal/pkg"
	"github.com/prometheus/client_golang/prometheus"

	forum_usecase "github.com/perlinleo/technopark-mail.ru-forum-database/internal/app/forum/usecase"
	thread_usecase "github.com/perlinleo/technopark-mail.ru-forum-database/internal/app/thread/usecase"
	user_usecase "github.com/perlinleo/technopark-mail.ru-forum-database/internal/app/user/usecase"
	"github.com/valyala/fasthttp"

	forum_http "github.com/perlinleo/technopark-mail.ru-forum-database/internal/app/forum/delivery"
	thread_http "github.com/perlinleo/technopark-mail.ru-forum-database/internal/app/thread/delivery"
	user_http "github.com/perlinleo/technopark-mail.ru-forum-database/internal/app/user/delivery"
)

func Start() error {

	var metrics responses.PromMetrics

	metrics.Hits = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "hits",
	}, []string{"status", "path", "method"})

	metrics.Timings = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "timings",
		},
		[]string{"status", "path", "method"},
	)
	*metrics.Requests = 0;
	
	prometheus.MustRegister(metrics.Hits, metrics.Timings)
	
	config := NewConfig()

	_, err := NewServer(config)
	if err != nil {
		return err
	}
	ConnPool, err := NewDataBase(config.App.DatabaseURL)
	if err != nil {
		return err
	}

	userCache := cache.New(5*time.Minute, 10*time.Minute)
	userRepository := user_psql.NewUserPSQLRepository(ConnPool, userCache)
	
	forumRepository := forum_psql.NewForumPSQLRepository(ConnPool, userCache)
	
	threadRepository := thread_psql.NewThreadPSQLRepository(ConnPool, userCache)

	threadUsecase := thread_usecase.NewThreadUsecase(threadRepository, userRepository)
	userUsecase := user_usecase.NewUserUsecase(userRepository)
	forumUsecase := forum_usecase.NewForumUsecase(forumRepository, threadRepository, userRepository, userCache)
	router := router.New()
	user_http.NewUserHandler(router, userUsecase,&metrics)
	forum_http.NewForumHandler(router, forumUsecase,&metrics)
	thread_http.NewThreadHandler(router, threadUsecase,&metrics)
	
	
	fmt.Printf("STARTING SERVICE ON PORT %s\n", config.App.Port)
	err = fasthttp.ListenAndServe(config.App.Port, router.Handler)
	if err != nil {
		return err
	}

	return nil
}
