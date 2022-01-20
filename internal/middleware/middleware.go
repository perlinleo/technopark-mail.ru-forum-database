package middleware

import (
	"github.com/valyala/fasthttp"
)

func ReponseMiddlwareAndLogger(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		ctx.Response.Header.Set("Content-Type", "application/json")
		// start := time.Now()
		// log.Printf("METHOD %s REMOTEADDR %s URL %s", ctx.Method(), ctx.RemoteAddr(), ctx.RequestURI())
		// elapsed := time.Since(start)
		// log.Printf("Time spent %s", elapsed)
		next(ctx)
	}
}
