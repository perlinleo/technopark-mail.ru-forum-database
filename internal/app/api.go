package app

import (
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/valyala/fasthttp"
)
func Start() error {
	config := NewConfig()

	server, err := NewServer(config)
	if err !=nil {
		return err
	}
	connPool, err := NewDataBase(config.App.DatabaseURL)
	if err !=nil {
		return err
	}

	fmt.Println(connPool)

	log.Error().Msgf(fasthttp.ListenAndServe(config.App.Port,server.Router.Handler).Error())
	
	return nil;
}