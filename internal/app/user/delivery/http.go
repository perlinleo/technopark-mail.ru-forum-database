package user_http

import (
	"encoding/json"
	"net/http"

	"github.com/fasthttp/router"
	"github.com/perlinleo/technopark-mail.ru-forum-database/internal/app/user"
	"github.com/perlinleo/technopark-mail.ru-forum-database/internal/middleware"
	"github.com/perlinleo/technopark-mail.ru-forum-database/internal/model"
	responses "github.com/perlinleo/technopark-mail.ru-forum-database/internal/pkg"
	"github.com/valyala/fasthttp"
)

type UserHandler struct {
	UserUsecase user.Usecase
}

func NewUserHandler(router *router.Router, usecase user.Usecase) {
	handler := &UserHandler{
		UserUsecase: usecase,
	}

	router.POST("/api/user/{nickname}/create", middleware.ReponseMiddlwareAndLogger(handler.HandleCreateUser))
	router.GET("/api/user/{nickname}/profile", middleware.ReponseMiddlwareAndLogger(handler.HandleGetUser))
	router.POST("/api/user/{nickname}/profile", middleware.ReponseMiddlwareAndLogger(handler.HandleUpdateUser))

}

func (h *UserHandler) HandleCreateUser(ctx *fasthttp.RequestCtx) {

	nickname := ctx.UserValue("nickname").(string)

	newUser := new(model.User)
	newUser.Nickname = nickname

	err := json.Unmarshal(ctx.PostBody(), &newUser)
	if err != nil {
		responses.SendError(ctx, err, fasthttp.StatusInternalServerError)
	}

	users, err := h.UserUsecase.CreateUser(newUser)
	if err != nil {
		// pgerr, _ := err.(pgx.PgError);
		users, err = h.UserUsecase.DuplicateUser(newUser)
		
		ctxBody, err := json.Marshal(users)
		if err != nil {
			responses.SendError(ctx, err, fasthttp.StatusConflict)
		}

		ctx.SetBody(ctxBody)
		ctx.SetStatusCode(fasthttp.StatusConflict)

		return
	}
	if users != nil {
		ctxBody, err := json.Marshal(users)
		if err != nil {
			responses.SendError(ctx, err, fasthttp.StatusInternalServerError)
		}
		ctx.SetBody(ctxBody)
		return
	}
	ctxBody, err := json.Marshal(newUser)
	ctx.SetBody(ctxBody)
	if err != nil {
		responses.SendError(ctx, err, fasthttp.StatusInternalServerError)
	}
	ctx.SetStatusCode(fasthttp.StatusCreated)

}

func (h *UserHandler) HandleGetUser(ctx *fasthttp.RequestCtx) {
	nickname := ctx.UserValue("nickname").(string)
	userObj, err := h.UserUsecase.Find(nickname)
	ctx.SetStatusCode(http.StatusOK)
	if err != nil || userObj == nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
	}

	userBody, err := json.Marshal(userObj)
	ctx.SetBody(userBody)
}

func (h *UserHandler) HandleUpdateUser(ctx *fasthttp.RequestCtx) {

	nickname := ctx.UserValue("nickname").(string)
	newUser := new(model.User)

	newUser.Nickname = nickname
	err := json.Unmarshal(ctx.PostBody(), &newUser)
	if err != nil {
		responses.SendError(ctx, err, fasthttp.StatusInternalServerError)
	}
	newUser, err, code := h.UserUsecase.Update(newUser)
	if err != nil {
		responses.SendError(ctx, err, fasthttp.StatusConflict)
	}
	ctxBody, err := json.Marshal(newUser)
	if err != nil {
		responses.SendError(ctx, err, fasthttp.StatusInternalServerError)
	}
	ctx.SetBody(ctxBody)
	ctx.SetStatusCode(code)

	// newUser.Nickname = nickname
	// newUser, err, code := h.UserUsecase.Update(newUser)
}
