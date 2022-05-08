package app

import (
	"fmt"
	"time"

	fasthttpprom "github.com/carousell/fasthttp-prometheus-middleware"
	"github.com/prometheus/client_golang/prometheus"

	router "github.com/fasthttp/router"
	"github.com/patrickmn/go-cache"
	forum_psql "github.com/perlinleo/technopark-mail.ru-forum-database/internal/app/forum/repository"
	thread_psql "github.com/perlinleo/technopark-mail.ru-forum-database/internal/app/thread/repository"
	user_psql "github.com/perlinleo/technopark-mail.ru-forum-database/internal/app/user/repository"

	forum_usecase "github.com/perlinleo/technopark-mail.ru-forum-database/internal/app/forum/usecase"
	thread_usecase "github.com/perlinleo/technopark-mail.ru-forum-database/internal/app/thread/usecase"
	user_usecase "github.com/perlinleo/technopark-mail.ru-forum-database/internal/app/user/usecase"
	"github.com/perlinleo/technopark-mail.ru-forum-database/internal/middleware"
	"github.com/valyala/fasthttp"

	responses "github.com/perlinleo/technopark-mail.ru-forum-database/internal/pkg"

	forum_http "github.com/perlinleo/technopark-mail.ru-forum-database/internal/app/forum/delivery"
	thread_http "github.com/perlinleo/technopark-mail.ru-forum-database/internal/app/thread/delivery"
	user_http "github.com/perlinleo/technopark-mail.ru-forum-database/internal/app/user/delivery"
)

func Index(ctx *fasthttp.RequestCtx) {
	fmt.Println("sdss")
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusOK)
}

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
	router := router.New()
	userCache := cache.New(time.Minute, time.Minute)
	userRepository := user_psql.NewUserPSQLRepository(ConnPool, userCache)
	forumCache := cache.New(time.Minute, time.Minute)
	forumRepository := forum_psql.NewForumPSQLRepository(ConnPool, forumCache)
	threadCache := cache.New(time.Minute, time.Minute)
	threadRepository := thread_psql.NewThreadPSQLRepository(ConnPool, threadCache)
	forumUsecaseCache := cache.New(time.Minute, time.Minute)
	router.GET("/",middleware.ReponseMiddlwareAndLogger(Index,&metrics))
	router.GET("/api/",middleware.ReponseMiddlwareAndLogger(Index, &metrics))
	threadUsecase := thread_usecase.NewThreadUsecase(threadRepository, userRepository)
	userUsecase := user_usecase.NewUserUsecase(userRepository)
	forumUsecase := forum_usecase.NewForumUsecase(forumRepository, threadRepository, userRepository, forumUsecaseCache)

	user_http.NewUserHandler(router, userUsecase,&metrics)
	forum_http.NewForumHandler(router, forumUsecase,&metrics)
	thread_http.NewThreadHandler(router, threadUsecase,&metrics)
	
	
	fmt.Printf("STARTING SERVICE ON PORT %s\n", config.App.Port)
	p := fasthttpprom.NewPrometheus("")
	
	p.Use(router)

	err = fasthttp.ListenAndServe(config.App.Port, middleware.ReponseMiddlwareAndLogger(router.Handler, &metrics))
	if err != nil {
		return err
	}

	return nil
}
