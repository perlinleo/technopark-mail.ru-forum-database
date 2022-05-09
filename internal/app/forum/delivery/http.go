package forum_http

import (
	"encoding/json"

	"github.com/fasthttp/router"
	"github.com/perlinleo/technopark-mail.ru-forum-database/internal/app/forum"
	"github.com/perlinleo/technopark-mail.ru-forum-database/internal/middleware"
	"github.com/perlinleo/technopark-mail.ru-forum-database/internal/model"
	responses "github.com/perlinleo/technopark-mail.ru-forum-database/internal/pkg"
	"github.com/valyala/fasthttp"
)

type ForumHandler struct {
	ForumUsecase forum.Usecase
}

func NewForumHandler(router *router.Router, usecase forum.Usecase,metrics *responses.PromMetrics) {
	handler := &ForumHandler{
		ForumUsecase: usecase,
	}
	router.POST("/api/forum/create", middleware.ReponseMiddlwareAndLogger(handler.CreateForum,metrics))
	router.POST("/api/forum/{slug}/create", middleware.ReponseMiddlwareAndLogger(handler.CreateThread,metrics))
	router.GET("/api/forum/{slug}/details", middleware.ReponseMiddlwareAndLogger(handler.GetDetails,metrics))
	router.GET("/api/forum/{slug}/threads", middleware.ReponseMiddlwareAndLogger(handler.GetThreads,metrics))
	router.GET("/api/forum/{slug}/users", middleware.ReponseMiddlwareAndLogger(handler.GetUsers,metrics))
	router.GET("/metrics", middleware.ReponseMiddlwareAndLogger(handler.Metrics,metrics))
}


func (h *ForumHandler) Metrics(ctx *fasthttp.RequestCtx) {
	ctxBody := []byte("dokjasdoiassadjkaso ")
	ctx.SetBody(ctxBody)
	ctx.SetStatusCode(200)
}


func (h *ForumHandler) CreateForum(ctx *fasthttp.RequestCtx) {
	newForum := new(model.Forum)
	err := json.Unmarshal(ctx.PostBody(), &newForum)
	if err != nil {
		responses.SendError(ctx, err, fasthttp.StatusInternalServerError)
		return
	}

	forumObj, code, err := h.ForumUsecase.CreateForum(newForum)
	if err != nil {
		responses.SendError(ctx, err, fasthttp.StatusInternalServerError)
	}
	ctxBody, err := json.Marshal(forumObj)
	ctx.SetBody(ctxBody)
	if code == fasthttp.StatusNotFound {
		responses.SendError(ctx, err, fasthttp.StatusNotFound)
		return
	}

	if code == fasthttp.StatusConflict {
		responses.SendError(ctx, err, fasthttp.StatusConflict)
		return
	}

	if err != nil {
		responses.SendError(ctx, err, fasthttp.StatusInternalServerError)
	}
	ctx.SetBody(ctxBody)
	ctx.SetStatusCode(fasthttp.StatusCreated)
}

func (h *ForumHandler) CreateThread(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("slug").(string)
	newThread := new(model.NewThread)
	err := json.Unmarshal(ctx.PostBody(), &newThread)
	if err != nil {
		responses.SendError(ctx, err, fasthttp.StatusBadRequest)
		return
	}
	ctx.SetBody(ctx.PostBody())
	threadObj, code, err := h.ForumUsecase.CreateThread(slug, newThread)

	ctxBody, err := json.Marshal(threadObj)
	if err == nil {
		ctx.SetBody(ctxBody)
	}
	
	ctx.SetStatusCode(code)
}

func (h *ForumHandler) GetDetails(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("slug").(string)
	forumObj, err := h.ForumUsecase.Find(slug)
	if err != nil || forumObj == nil {
		responses.SendError(ctx, err, fasthttp.StatusNotFound)
		response := map[string]string{"message": "Can`t find threads for forum slug"}
		ctxBody, _ := json.Marshal(response)
		ctx.SetBody(ctxBody)
		ctx.SetStatusCode(fasthttp.StatusNotFound)

		return
	}
	ctxBody, err := json.Marshal(forumObj)
	ctx.SetBody(ctxBody)
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func (h *ForumHandler) GetThreads(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("slug").(string)
	args := ctx.QueryArgs()

	var limitValue string
	var descValue bool
	var sinceValue string
	limitValue = string(args.Peek("limit"))
	if string(args.Peek("desc")) == "true" {
		descValue = true
	}
	sinceValue = string(args.Peek("since"))

	threads, code, err := h.ForumUsecase.GetThreadsByForum(slug, limitValue, descValue, sinceValue)
	ctxBody, err := json.Marshal(threads)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBody([]byte(err.Error()))
	}
	ctx.SetBody(ctxBody)
	ctx.SetStatusCode(code)

}

func (h *ForumHandler) GetUsers(ctx *fasthttp.RequestCtx) {

	slug := ctx.UserValue("slug").(string)
	args := ctx.QueryArgs()

	var limitValue string
	var descValue bool
	var sinceValue string
	limitValue = string(args.Peek("limit"))
	if string(args.Peek("desc")) == "true" {
		descValue = true
	}
	sinceValue = string(args.Peek("since"))

	users, code, err := h.ForumUsecase.GetUsersByForum(slug, limitValue, descValue, sinceValue)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
	}
	ctxBody, err := json.Marshal(users)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBody([]byte(err.Error()))
	}
	ctx.SetBody(ctxBody)
	ctx.SetStatusCode(code)
}
