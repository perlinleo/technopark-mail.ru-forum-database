package thread_http

import (
	"encoding/json"
	"fmt"

	"github.com/fasthttp/router"
	"github.com/perlinleo/technopark-mail.ru-forum-database/internal/app/thread"
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

	router.POST("/thread/{slug}/create", handler.HandleCreatePosts)
	router.GET("/thread/{slug}/details", handler.HandleGetThreadDetails)
	router.POST("/thread/{slug}/details", handler.HandleUpdateThreadDetails)
	router.GET("/thread/{slug}/posts", handler.HandleGetThreadPosts)
	router.POST("/thread/{slug}/vote", handler.HandleVoteForThread)
}

func (h *ThreadHandler) HandleCreatePosts(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("slug").(string)

	newPosts := new([]*model.Post)

	err := json.Unmarshal(ctx.PostBody(), &newPosts)
	
	if err != nil {
		fmt.Println(err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	posts, code, err := h.ThreadUsecase.CreatePosts(slug, *newPosts)
	
	if err != nil {
		fmt.Println(err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(code)
	
	ctxBody , err:= json.Marshal(posts)

	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}
	
	ctx.SetBody(ctxBody)

}

func (h *ThreadHandler) HandleGetThreadDetails(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("slug").(string)

	threadObj, err := h.ThreadUsecase.FindByIdOrSlug(slug)

	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
	}

	ctxBody , err:= json.Marshal(threadObj)

	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
	}
	
	ctx.SetBody(ctxBody)
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func (h *ThreadHandler) HandleUpdateThreadDetails(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("slug").(string)

	thread := new(model.ThreadUpdate)

	err := json.Unmarshal(ctx.PostBody(), &thread)
	
	if err != nil {
		fmt.Println(err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
	threadObj, err := h.ThreadUsecase.UpdateThread(slug, thread)
	if err != nil {
		fmt.Println(err)
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}

	ctxBody , err:= json.Marshal(threadObj)

	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
	
	ctx.SetBody(ctxBody)
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func (h *ThreadHandler) HandleGetThreadPosts(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("slug").(string)
	args := ctx.QueryArgs()
	
	var limitValue string
	var descValue bool
	var sinceValue string
	var sortValue string

	limitValue = string(args.Peek("limit"))
	if string(args.Peek("desc")) == "true"{
		descValue = true
	}
	sinceValue = string(args.Peek("since"))
	sortValue = string(args.Peek("sort"))
	posts, err := h.ThreadUsecase.GetThreadPosts(slug, limitValue, descValue, sinceValue, sortValue)


	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	ctxBody , err:= json.Marshal(posts)

	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
	
	ctx.SetBody(ctxBody)
	ctx.SetStatusCode(fasthttp.StatusOK)

}

func (h *ThreadHandler) HandleVoteForThread(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("slug").(string)
	vote := new(model.Vote)

	err := json.Unmarshal(ctx.PostBody(), &vote)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
	}


	threadObj ,err := h.ThreadUsecase.Vote(slug, vote)
	if err != nil {
		fmt.Println(err)
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
