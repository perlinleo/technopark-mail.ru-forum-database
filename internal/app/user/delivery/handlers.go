package delivery

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fasthttp/router"
	"github.com/jackc/pgx"
	"github.com/perlinleo/technopark-mail.ru-forum-database/internal/app/user"
	"github.com/perlinleo/technopark-mail.ru-forum-database/internal/middleware"
	"github.com/perlinleo/technopark-mail.ru-forum-database/internal/model"
	responses "github.com/perlinleo/technopark-mail.ru-forum-database/internal/pkg"
	"github.com/valyala/fasthttp"
)

type UserHandler struct {
	UserUsecase user.Usecase
}

func NewUserHandler(router *router.Router,usecase user.Usecase) {
	handler := &UserHandler{
		UserUsecase: usecase,
	}

	router.POST("/user/{nickname}/create", middleware.ReponseMiddlwareAndLogger(handler.HandleCreateUser))
	router.GET("/user/{nickname}/profile", middleware.ReponseMiddlwareAndLogger(handler.HandleGetUser))
	router.POST("/user/{nickname}/profile", middleware.ReponseMiddlwareAndLogger(handler.HandleUpdateUser))
	
}


func (h *UserHandler) HandleCreateUser(ctx *fasthttp.RequestCtx) {
	

	nickname := ctx.UserValue("nickname").(string)
	fmt.Println(nickname)

	newUser := new(model.User)
	newUser.Nickname = nickname

	err := json.Unmarshal(ctx.PostBody(), &newUser)
	if err != nil {
		responses.SendError(ctx,err, fasthttp.StatusInternalServerError)
	}

	users, err := h.UserUsecase.CreateUser(newUser)
	if err != nil{
		pgerr, _ := err.(pgx.PgError);
		fmt.Println(pgerr.ConstraintName)

			fmt.Println("POYMALI DOLBAEBA")
			users,err = h.UserUsecase.DuplicateUser(newUser)
			ctxBody, err := json.Marshal(users)
			if err != nil{
				fmt.Println("pohuy")
			}
			ctx.SetBody(ctxBody)
			ctx.SetStatusCode(fasthttp.StatusConflict)
		return
	}
	fmt.Println(users)
	if users!=nil {
		ctxBody, err := json.Marshal(users)
		fmt.Println(users)
		if err != nil {
			responses.SendError(ctx,err, fasthttp.StatusInternalServerError)
		}
		ctx.SetBody(ctxBody)
		return
	}
	ctxBody, err := json.Marshal(newUser)
	ctx.SetBody(ctxBody)
	if err != nil {
		responses.SendError(ctx,err, fasthttp.StatusInternalServerError)
	}
	fmt.Println("hello")
	ctx.SetStatusCode(fasthttp.StatusCreated)
	
}

func (h *UserHandler) HandleGetUser(ctx *fasthttp.RequestCtx) {
	nickname := ctx.UserValue("nickname").(string)
	userObj, err := h.UserUsecase.Find(nickname)

	if err != nil || userObj == nil {
		responses.SendError(ctx,err, fasthttp.StatusNotFound)
	}
	ctx.SetStatusCode(http.StatusOK)
	userBody,err := json.Marshal(userObj);
	ctx.SetBody(userBody)
}

func (h *UserHandler) HandleUpdateUser(ctx *fasthttp.RequestCtx) {
	nickname := ctx.UserValue("nickname").(string)
	newUser := new(model.User)
	newUser.Nickname = nickname
	err := json.Unmarshal(ctx.PostBody(), &newUser)
	if err != nil {
		fmt.Println(err)
	}
	newUser, err , _ = h.UserUsecase.Update(newUser)
	if err != nil {
		fmt.Println(err)
	}
	ctxBody, err := json.Marshal(newUser)
	ctx.SetBody(ctxBody)
	// newUser.Nickname = nickname
	// newUser, err, code := h.UserUsecase.Update(newUser)
} 