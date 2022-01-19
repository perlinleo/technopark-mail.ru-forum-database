package app

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type DataBaseConfig struct {
	host string
	port string
	user string
	pass string
	name string
}

type AppConfig struct {
	SessionKey  string
	Port 	    string
	DatabaseURL string
}

type Config struct {
	App *AppConfig
	DB  *DataBaseConfig 
}

func NewConfig() *Config {
	viper.SetConfigFile("/home/pierrelean/technopark-mail.ru-forum-database/config.json")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	if viper.GetBool(`debug`) {
		log.Println("Service RUN on DEBUG mode")
	}

	dbHost := viper.GetString(`database.host`)
	dbPort := viper.GetString(`database.port`)
	dbUser := viper.GetString(`database.user`)
	dbPass := viper.GetString(`database.pass`)
	dbName := viper.GetString(`database.name`)

	appPort := viper.GetString(`server.address`)

	connection := fmt.Sprintf(
		"host=%s dbname=%s sslmode=disable port=%s password=%s user=%s",
		 dbHost, dbName, dbPort, dbPass, dbUser)
	return &Config{
		DB : &DataBaseConfig{
			host: dbHost,
			port: dbPort,
			user: dbUser,
			pass: dbPass,
			name: dbName,
		},
		App : &AppConfig{
			SessionKey: "dsdsdsds",
			DatabaseURL: connection,
			Port: appPort,
		},
	}
}