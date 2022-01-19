package responses

import (
	"encoding/json"

	"github.com/valyala/fasthttp"
)


func SendError(ctx *fasthttp.RequestCtx,errorMessage error, statusCode int) {
	ctx.SetStatusCode(statusCode)
}

func SendResponse(data interface{},ctx *fasthttp.RequestCtx,StatusCode int) {
	ctx.SetStatusCode(StatusCode)

	serializedData, err := json.Marshal(data)
	if err != nil {
		SendError(ctx,err,500)
	}

	ctx.SetBody(serializedData)
}
