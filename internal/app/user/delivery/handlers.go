package delivery

import (
	"encoding/json"

	"github.com/fasthttp/router"
	"github.com/perlinleo/technopark-mail.ru-forum-database/internal/app/user"
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

	router.POST("/user/{nickname}/create", handler.HandleCreateUser)
	// router.GET("/user/{nickname}/profile", handler.HandleCreateUser)
	// router.POST("/user/{nickname}/profile", handler.HandleCreateUser)
	
}


func (h *UserHandler) HandleCreateUser(ctx *fasthttp.RequestCtx) {
	

	nickname := ctx.QueryArgs().Peek("nickname")

	newUser := model.User{Nickname: string(nickname)}

	err := json.Unmarshal(ctx.PostBody(), &newUser)

	if err != nil {
		responses.SendError(err,ctx)
	}

	err = h.UserUsecase.CreateUser(&newUser)

	if err != nil {
		responses.SendError(err,ctx)
	}

}