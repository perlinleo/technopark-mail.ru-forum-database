package forumApp

import (
	"fmt"

	router "github.com/fasthttp/router"
	"github.com/jackc/pgx"
	"github.com/valyala/fasthttp"
)

type Server struct {
	Router *router.Router
	Config *Config
}

func Index(ctx *fasthttp.RequestCtx){
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func NewServer(config *Config) (*Server, error) {
	server:= &Server{
		Router: NewRouter(),
		Config: config,
	}

	return server, nil
}
func NewRouter() (*router.Router) {
	router:= router.New()
	router.GET("/", Index)

	return router;
}

func NewDataBase(connectionString string) (*pgx.ConnPool, error){
	pgxConn, err := pgx.ParseConnectionString(connectionString)
	if err != nil{
		return nil, err
	}
	pgxConn.PreferSimpleProtocol = true
	
	config := pgx.ConnPoolConfig{
		ConnConfig:     pgxConn,
		MaxConnections: 100,
		AfterConnect:   nil,
		AcquireTimeout: 0,
	}


	connPool, err := pgx.NewConnPool(config)
	if err != nil{
		return nil, err
	}

	return connPool, nil

}


func (s *Server) ConfigurateServer(pgx.ConnPool) {
	fmt.Println("configurating")


}