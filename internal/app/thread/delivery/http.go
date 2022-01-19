package thread_http

import (
	"encoding/json"
	"fmt"
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

	router.POST("/thread/{slug}/create",  middleware.ReponseMiddlwareAndLogger(handler.HandleCreatePosts))
	router.GET("/thread/{slug}/details",  middleware.ReponseMiddlwareAndLogger(handler.HandleGetThreadDetails))
	router.POST("/thread/{slug}/details", middleware.ReponseMiddlwareAndLogger(handler.HandleUpdateThreadDetails))
	router.GET("/thread/{slug}/posts",  middleware.ReponseMiddlwareAndLogger(handler.HandleGetThreadPosts))
	router.POST("/thread/{slug}/vote",  middleware.ReponseMiddlwareAndLogger(handler.HandleVoteForThread))
	router.GET("/post/{id}/details",  middleware.ReponseMiddlwareAndLogger(handler.HandleGetPostDetails))
	router.POST("/post/{id}/details",  middleware.ReponseMiddlwareAndLogger(handler.HandleUpdatePost))

	router.POST("/service/clear", middleware.ReponseMiddlwareAndLogger(handler.HandleServiceClear))
	router.GET("/service/status", middleware.ReponseMiddlwareAndLogger(handler.HandleServiceGetStatus))
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
		// fmt.Println(err)
		// ctx.SetStatusCode(fasthttp.StatusInternalServerError)

		// if string(err.Error()) == "Parent post was created in another thread" {
			
		// }
		fmt.Println(err)
		if string(err.Error()) == "404" {
			response := map[string]string{"message": "Can't find post author by nickname: "}
			ctxBody, _ := json.Marshal(response)
			ctx.SetBody(ctxBody)
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			return
		}
		if string(err.Error()) == "no rows in result set" && code!= 409 {
			fmt.Println("???")
			fmt.Println(code)
			response := map[string]string{"message": "Can't find post thread by id: 2139800939"}
			ctxBody, _ := json.Marshal(response)
			ctx.SetBody(ctxBody)
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			return
		}
		fmt.Println("!!!")
		response := map[string]string{"message": "Parent post was created in another thread"}
		ctxBody, _ := json.Marshal(response)
		ctx.SetBody(ctxBody)
		ctx.SetStatusCode(fasthttp.StatusConflict)
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
		fmt.Println(err)
		// mb kringe
		response := map[string]string{"message": "Can't find user by nickname:  "}
		ctxBody, _ := json.Marshal(response)
		ctx.SetBody(ctxBody)
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
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
		// mb kringe
		fmt.Println(err)
		response := map[string]string{"message": "Can't find user by nickname:  "}
		ctxBody, _ := json.Marshal(response)
		ctx.SetBody(ctxBody)
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
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		response := map[string]string{"message": "Can't find user by nickname:  "}
		ctxBody, _ := json.Marshal(response)
		ctx.SetBody(ctxBody)
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

func (h *ThreadHandler) HandleUpdatePost(ctx *fasthttp.RequestCtx) {
	id := ctx.UserValue("id").(string)
	newPost := new(model.Post)
	err := json.Unmarshal(ctx.PostBody(), &newPost)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	postObj, err := h.ThreadUsecase.UpdatePost(id, newPost.Message)

	if err != nil || postObj == nil {
		// respond.Error(w, r, http.StatusNotFound, errors.New("Can't find post with id "+id+"\n"))
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

func (h *ThreadHandler) HandleGetPostDetails(ctx *fasthttp.RequestCtx) {
	id := ctx.UserValue("id").(string)
	args := ctx.QueryArgs()
	fmt.Println(args)
	fmt.Println(id)
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
	postObj, err := h.ThreadUsecase.FindPostId(id,includeUser,includeForum,includeThread)
	fmt.Println(err)
	fmt.Println(postObj)
	// fmt.Println(postObj.Author == nil)
	fmt.Println("___________")
	if err != nil || postObj == nil  {
		// respond.Error(w, r, http.StatusNotFound, errors.New("Can't find post with id "+id+"\n"))
		// mb kringe
		fmt.Println(err)
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

	// !!!!!
	// Я ХОЧУ УМРЕРЕТЬ
	// ПЕТЯ Я СЕЙЧАС УМРУ
	// ПЕТЯ Я СЕЙЧАС СДОХНУ

	// respond.Respond(w, r, http.StatusOK, postObj)
}

func (h *ThreadHandler) HandleServiceClear(ctx *fasthttp.RequestCtx) {

	err := h.ThreadUsecase.ClearAll()

	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
}

func (h *ThreadHandler) HandleServiceGetStatus(ctx *fasthttp.RequestCtx) {

	status, err := h.ThreadUsecase.GetStatus()

	if err != nil || status == nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}
	ctxBody,err:= json.Marshal(status)
	ctx.SetBody(ctxBody)
	ctx.SetStatusCode(fasthttp.StatusOK)
}