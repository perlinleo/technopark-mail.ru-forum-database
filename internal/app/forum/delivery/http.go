package delivery

import (
	"github.com/fasthttp/router"
	"github.com/perlinleo/technopark-mail.ru-forum-database/internal/app/forum"
	"github.com/perlinleo/technopark-mail.ru-forum-database/internal/middleware"
	"github.com/valyala/fasthttp"
)

type ForumHandler struct {
	ForumUsecase forum.Usecase
}

func NewForumHandler(router *router.Router,usecase forum.Usecase) {
	handler := &ForumHandler{
		ForumUsecase: usecase,
	}

	router.POST("/forum/create", middleware.ReponseMiddlwareAndLogger(handler.CreateForum))
	router.POST("/forum/{slug}/create", middleware.ReponseMiddlwareAndLogger(handler.CreateThread))
	router.GET("/forum/{slug}/details", middleware.ReponseMiddlwareAndLogger(handler.GetDetails))
	router.GET("/forum/{slug}/threads", middleware.ReponseMiddlwareAndLogger(handler.GetThreads))
	router.GET("/forum/{slug}/users", middleware.ReponseMiddlwareAndLogger(handler.GetUsers))
}

func (h *ForumHandler) CreateForum(ctx *fasthttp.RequestCtx) {
	
}


func (h *ForumHandler) CreateThread(ctx *fasthttp.RequestCtx) {
	
}


func (h *ForumHandler) GetDetails(ctx *fasthttp.RequestCtx) {
	
}


func (h *ForumHandler) GetThreads(ctx *fasthttp.RequestCtx) {
	
}


func (h *ForumHandler) GetUsers(ctx *fasthttp.RequestCtx) {
	
}