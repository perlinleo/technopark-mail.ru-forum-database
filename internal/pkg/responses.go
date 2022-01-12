package responses

import (
	"encoding/json"
	"fmt"

	"github.com/valyala/fasthttp"
)


func SendError(errorMessage error, ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusInternalServerError)
	fmt.Printf("%s", errorMessage)
}

func SendResponse(data interface{},ctx *fasthttp.RequestCtx,StatusCode int) {
	ctx.SetStatusCode(StatusCode)

	serializedData, err := json.Marshal(data)
	if err != nil {
		SendError(err,ctx)
	}

	ctx.SetBody(serializedData)
}
