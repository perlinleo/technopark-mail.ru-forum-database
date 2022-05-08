package middleware

import (
	"net/http"
	"strconv"
	"time"

	responses "github.com/perlinleo/technopark-mail.ru-forum-database/internal/pkg"
	"github.com/valyala/fasthttp"
)

func ReponseMiddlwareAndLogger(next fasthttp.RequestHandler, metrics *responses.PromMetrics) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		ctx.Response.Header.Set("Content-Type", "application/json")
		// start := time.Now()
		// log.Printf("METHOD %s REMOTEADDR %s URL %s", ctx.Method(), ctx.RemoteAddr(), ctx.RequestURI())
		// elapsed := time.Since(start)
		// log.Printf("Time spent %s", elapsed)
		reqTime := time.Now()
		respTime := time.Since(reqTime)
		metrics.Hits.WithLabelValues(strconv.Itoa(http.StatusOK), string(ctx.RequestURI()), string(ctx.Method())).Inc()
		metrics.Timings.WithLabelValues(
					strconv.Itoa(http.StatusOK), string(ctx.RequestURI()), string(ctx.Method())).Observe(respTime.Seconds())
		next(ctx)
	}
}