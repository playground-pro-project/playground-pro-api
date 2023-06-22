package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

var JWT string

type AppConfig struct {
	DBUSER                string
	DBPASSWORD            string
	DBHOST                string
	DBPORT                string
	DBNAME                string
	AWS_ACCESS_KEY_ID     string
	AWS_SECRET_ACCESS_KEY string
	ADMINPASSWORD         string
}

func InitConfig() *AppConfig {
	return readEnv()
}

func readEnv() *AppConfig {
	app := AppConfig{}
	isRead := true

	if val, found := os.LookupEnv("DBUSER"); found {
		app.DBUSER = val
		isRead = false
	}

	if val, found := os.LookupEnv("DBPASSWORD"); found {
		app.DBPASSWORD = val
		isRead = false
	}

	if val, found := os.LookupEnv("DBHOST"); found {
		app.DBHOST = val
		isRead = false
	}

	if val, found := os.LookupEnv("DBPORT"); found {
		app.DBPORT = val
		isRead = false
	}

	if val, found := os.LookupEnv("DBNAME"); found {
		app.DBNAME = val
		isRead = false
	}

	if val, found := os.LookupEnv("JWT"); found {
		JWT = val
		isRead = false
	}

	if val, found := os.LookupEnv("ADMINPASSWORD"); found {
		app.ADMINPASSWORD = val
		isRead = false
	}

	if val, found := os.LookupEnv("AWS_ACCESS_KEY_ID"); found {
		app.AWS_ACCESS_KEY_ID = val
		isRead = false
	}

	if val, found := os.LookupEnv("AWS_SECRET_ACCESS_KEY"); found {
		app.AWS_SECRET_ACCESS_KEY = val
		isRead = false
	}

	if isRead {
		viper.AddConfigPath(".")
		viper.SetConfigName("local")
		viper.SetConfigType("yaml")

		err := viper.ReadInConfig()
		if err != nil {
			log.Println("error read config : ", err.Error())
			return nil
		}

		app.DBUSER = viper.GetString("DBUSER")
		app.DBPASSWORD = viper.GetString("DBPASSWORD")
		app.DBHOST = viper.GetString("DBHOST")
		app.DBPORT = viper.GetString("DBPORT")
		app.DBNAME = viper.GetString("DBNAME")
		JWT = viper.GetString("JWT")
		app.ADMINPASSWORD = viper.GetString("ADMINPASSWORD")
		app.AWS_ACCESS_KEY_ID = viper.Get("AWS_ACCESS_KEY_ID").(string)
		app.AWS_SECRET_ACCESS_KEY = viper.Get("AWS_SECRET_ACCESS_KEY").(string)
	}

	return &app
}
