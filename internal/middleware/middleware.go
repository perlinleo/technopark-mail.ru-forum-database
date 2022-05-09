package middleware

import (
	"log"
	"time"

	responses "github.com/perlinleo/technopark-mail.ru-forum-database/internal/pkg"
	"github.com/valyala/fasthttp"
)

func ReponseMiddlwareAndLogger(next fasthttp.RequestHandler, metrics *responses.PromMetrics) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		ctx.Response.Header.Set("Content-Type", "application/json")
		start := time.Now()
		log.Printf("METHOD %s REMOTEADDR %s URL %s", ctx.Method(), ctx.RemoteAddr(), ctx.RequestURI())
		elapsed := time.Since(start)
		log.Printf("Time spent %s", elapsed)
		ctx.SetUserValue("requests", metrics.Requests)
		metrics.Requests+=1;
		log.Printf("Requests total %s",string(rune(metrics.Requests)))
		next(ctx)
	}
}