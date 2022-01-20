package thread_http

import (
	"encoding/json"
	"strings"

	"github.com/fasthttp/router"
	"github.com/perlinleo/technopark-mail.ru-forum-database/internal/app/thread"
	"github.com/perlinleo/technopark-mail.ru-forum-database/internal/middleware"
	"github.com/perlinleo/technopark-mail.ru-forum-database/internal/model"
	"github.com/valyala/fasthttp"
)

type ThreadHandler struct {
	ThreadUsecase thread.Usecase
}

func NewThreadHandler(router *router.Router, u thread.Usecase) {
	handler := &ThreadHandler{
		ThreadUsecase: u,
	}

	router.POST("/api/thread/{slug}/create", middleware.ReponseMiddlwareAndLogger(handler.CreatePosts))
	router.GET("/api/thread/{slug}/details", middleware.ReponseMiddlwareAndLogger(handler.GetThreadDetails))
	router.POST("/api/thread/{slug}/details", middleware.ReponseMiddlwareAndLogger(handler.UpdateThreadDetails))
	router.GET("/api/thread/{slug}/posts", middleware.ReponseMiddlwareAndLogger(handler.GetThreadPosts))
	router.POST("/api/thread/{slug}/vote", middleware.ReponseMiddlwareAndLogger(handler.VoteForThread))
	router.GET("/api/post/{id}/details", middleware.ReponseMiddlwareAndLogger(handler.GetPostDetails))
	router.POST("/api/post/{id}/details", middleware.ReponseMiddlwareAndLogger(handler.UpdatePost))

	router.POST("/api/service/clear", middleware.ReponseMiddlwareAndLogger(handler.ServiceClear))
	router.GET("/api/service/status", middleware.ReponseMiddlwareAndLogger(handler.ServiceGetStatus))
}

func (h *ThreadHandler) CreatePosts(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("slug").(string)

	newPosts := new([]*model.Post)

	err := json.Unmarshal(ctx.PostBody(), &newPosts)

	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	posts, code, err := h.ThreadUsecase.CreatePosts(slug, *newPosts)

	if err != nil {
		
		if string(err.Error()) == "404" {
			response := map[string]string{"message": "Can't find post author by nickname: "}
			ctxBody, _ := json.Marshal(response)
			ctx.SetBody(ctxBody)
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			return
		}
		if string(err.Error()) == "no rows in result set" && code != 409 {
			response := map[string]string{"message": "Can't find post thread by id: 2139800939"}
			ctxBody, _ := json.Marshal(response)
			ctx.SetBody(ctxBody)
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			return
		}
		response := map[string]string{"message": "Parent post was created in another thread"}
		ctxBody, _ := json.Marshal(response)
		ctx.SetBody(ctxBody)
		ctx.SetStatusCode(fasthttp.StatusConflict)
		return
	}

	ctx.SetStatusCode(code)

	ctxBody, err := json.Marshal(posts)

	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	ctx.SetBody(ctxBody)

}

func (h *ThreadHandler) GetThreadDetails(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("slug").(string)

	threadObj, err := h.ThreadUsecase.FindByIdOrSlug(slug)

	if err != nil {
		response := map[string]string{"message": "Can't find user by nickname:  "}
		ctxBody, _ := json.Marshal(response)
		ctx.SetBody(ctxBody)
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}

	ctxBody, err := json.Marshal(threadObj)

	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
	}

	ctx.SetBody(ctxBody)
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func (h *ThreadHandler) UpdateThreadDetails(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("slug").(string)

	thread := new(model.ThreadUpdate)

	err := json.Unmarshal(ctx.PostBody(), &thread)

	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
	threadObj, err := h.ThreadUsecase.UpdateThread(slug, thread)
	if err != nil {
		response := map[string]string{"message": "Can't find user by nickname:  "}
		ctxBody, _ := json.Marshal(response)
		ctx.SetBody(ctxBody)
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}

	ctxBody, err := json.Marshal(threadObj)

	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetBody(ctxBody)
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func (h *ThreadHandler) GetThreadPosts(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("slug").(string)
	args := ctx.QueryArgs()

	var limitValue string
	var descValue bool
	var sinceValue string
	var sortValue string

	limitValue = string(args.Peek("limit"))
	if string(args.Peek("desc")) == "true" {
		descValue = true
	}
	sinceValue = string(args.Peek("since"))
	sortValue = string(args.Peek("sort"))
	posts, err := h.ThreadUsecase.GetThreadPosts(slug, limitValue, descValue, sinceValue, sortValue)

	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		response := map[string]string{"message": "Can't find user by nickname:  "}
		ctxBody, _ := json.Marshal(response)
		ctx.SetBody(ctxBody)
		return
	}

	ctxBody, err := json.Marshal(posts)

	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetBody(ctxBody)
	ctx.SetStatusCode(fasthttp.StatusOK)

}

func (h *ThreadHandler) VoteForThread(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("slug").(string)
	vote := new(model.Vote)

	err := json.Unmarshal(ctx.PostBody(), &vote)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
	}

	threadObj, err := h.ThreadUsecase.Vote(slug, vote)

	if err != nil {

		response := map[string]string{"message": "Can't find user by nickname:  "}
		ctxBody, _ := json.Marshal(response)
		ctx.SetBody(ctxBody)
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}
	ctxBody, err := json.Marshal(threadObj)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
	ctx.SetBody(ctxBody)
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func (h *ThreadHandler) UpdatePost(ctx *fasthttp.RequestCtx) {
	id := ctx.UserValue("id").(string)
	newPost := new(model.Post)
	err := json.Unmarshal(ctx.PostBody(), &newPost)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	postObj, err := h.ThreadUsecase.UpdatePost(id, newPost.Message)

	if err != nil || postObj == nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		response := map[string]string{"message": "Can't find post author by nickname: "}
		ctxBody, _ := json.Marshal(response)
		ctx.SetBody(ctxBody)
		return
	}

	ctxBody, err := json.Marshal(postObj)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
	ctx.SetBody(ctxBody)
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}

	return false
}

func (h *ThreadHandler) GetPostDetails(ctx *fasthttp.RequestCtx) {
	id := ctx.UserValue("id").(string)
	args := ctx.QueryArgs()
	var includeUser bool
	var includeForum bool
	var includeThread bool
	related := string(args.Peek("related"))
	if len(related) >= 1 {
		splitRelated := strings.Split(related, ",")

		if contains(splitRelated, "user") {
			includeUser = true
		}
		if contains(splitRelated, "forum") {
			includeForum = true
		}
		if contains(splitRelated, "thread") {
			includeThread = true
		}
	}
	postObj, err := h.ThreadUsecase.FindPostId(id, includeUser, includeForum, includeThread)
	if err != nil || postObj == nil {
		

		ctx.SetStatusCode(fasthttp.StatusNotFound)

		response := map[string]string{"message": "Can't find user by nickname:  "}
		ctxBody, _ := json.Marshal(response)
		ctx.SetBody(ctxBody)
		return
	}
	ctxBody, err := json.Marshal(postObj)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
	ctx.SetBody(ctxBody)
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func (h *ThreadHandler) ServiceClear(ctx *fasthttp.RequestCtx) {

	err := h.ThreadUsecase.ClearAll()

	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
}

func (h *ThreadHandler) ServiceGetStatus(ctx *fasthttp.RequestCtx) {

	status, err := h.ThreadUsecase.GetStatus()

	if err != nil || status == nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}
	ctxBody, err := json.Marshal(status)
	ctx.SetBody(ctxBody)
	ctx.SetStatusCode(fasthttp.StatusOK)
}
