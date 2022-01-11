package delivery

import (
	"../../user"
	"github.com/gorilla/mux"
	"github.com/valyala/fasthttp"
)

type UserHandler struct {
	UserUsecase user.Usecase
}

func UserHandler(m *mux.Router, u user.Usecase) {
	handler := &UserHandler{
		UserUsecase: u,
	}

	r.POST("/api/user/{nickname}/create", handler.Add)
}


func (ur *UserHandler) Add(ctx *fasthttp.RequestCtx) {
	nickname.
}