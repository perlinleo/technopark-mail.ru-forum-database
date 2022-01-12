package main

import (
	app "github.com/perlinleo/technopark-mail.ru-forum-database/internal/app"
	"github.com/rs/zerolog/log"
)


func main() {
	log.Error().Msgf(app.Start().Error());
}



func init() {
	
}