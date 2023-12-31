package config

import (
	"log"
	"os"
	"strconv"

	"github.com/spf13/viper"
)

var (
	err                   error
	JWT                   string
	REDIS_HOST            string
	REDIS_PORT            string
	REDIS_PASSWORD        string
	REDIS_DATABASE        int
	MIDTRANS_SERVERKEY    string
	MIDTRANS_MERCHANT_ID  string
	EMAIL_SENDER_NAME     string
	EMAIL_SENDER_ADDRESS  string
	EMAIL_SENDER_PASSWORD string
)

type AppConfig struct {
	DBUSER                string
	DBPASS                string
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

	if val, found := os.LookupEnv("DBPASS"); found {
		app.DBPASS = val
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

	if val, found := os.LookupEnv("REDIS_HOST"); found {
		REDIS_HOST = val
		isRead = false
	}

	if val, found := os.LookupEnv("REDIS_PORT"); found {
		REDIS_PORT = val
		isRead = false
	}

	if val, found := os.LookupEnv("REDIS_PASSWORD"); found {
		REDIS_PASSWORD = val
		isRead = false
	}

	if val, found := os.LookupEnv("REDIS_DATABASE"); found {
		REDIS_DATABASE, err = strconv.Atoi(val)
		if err != nil {
			log.Println("can't convert string to int")
		}
		isRead = false
	}

	if val, found := os.LookupEnv("MIDTRANS_SERVERKEY"); found {
		MIDTRANS_SERVERKEY = val
		isRead = false
	}

	if val, found := os.LookupEnv("MIDTRANS_MERCHANT_ID"); found {
		MIDTRANS_MERCHANT_ID = val
		isRead = false
	}

	if val, found := os.LookupEnv("EMAIL_SENDER_NAME"); found {
		EMAIL_SENDER_NAME = val
		isRead = false
	}

	if val, found := os.LookupEnv("EMAIL_SENDER_ADDRESS"); found {
		EMAIL_SENDER_ADDRESS = val
		isRead = false
	}

	if val, found := os.LookupEnv("EMAIL_SENDER_PASSWORD"); found {
		EMAIL_SENDER_PASSWORD = val
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

		JWT = viper.GetString("JWT")
		app.DBUSER = viper.GetString("DBUSER")
		app.DBPASS = viper.GetString("DBPASS")
		app.DBHOST = viper.GetString("DBHOST")
		app.DBPORT = viper.GetString("DBPORT")
		app.DBNAME = viper.GetString("DBNAME")
		app.ADMINPASSWORD = viper.GetString("ADMINPASSWORD")
		app.AWS_ACCESS_KEY_ID = viper.Get("AWS_ACCESS_KEY_ID").(string)
		app.AWS_SECRET_ACCESS_KEY = viper.Get("AWS_SECRET_ACCESS_KEY").(string)
		REDIS_HOST = viper.GetString("REDIS_HOST")
		REDIS_PORT = viper.GetString("REDIS_PORT")
		REDIS_PASSWORD = viper.GetString("REDIS_PASSWORD")
		REDIS_DATABASE = viper.GetInt("REDIS_DATABASE")
		MIDTRANS_SERVERKEY = viper.GetString("MIDTRANS_SERVERKEY")
		MIDTRANS_MERCHANT_ID = viper.GetString("MIDTRANS_MERCHANT_ID")
		EMAIL_SENDER_ADDRESS = viper.GetString("EMAIL_SENDER_ADDRESS")
		EMAIL_SENDER_NAME = viper.GetString("EMAIL_SENDER_NAME")
		EMAIL_SENDER_PASSWORD = viper.GetString("EMAIL_SENDER_PASSWORD")
	}

	return &app
}
