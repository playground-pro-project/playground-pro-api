package config

import (
	"log"
	"os"
	"strconv"

	"github.com/spf13/viper"
)

var (
	err                error
	JWT                string
	OTP                string
	REDIS_HOST         string
	REDIS_PORT         string
	REDIS_PASSWORD     string
	REDIS_DATABASE     int
	MIDTRANS_SERVERKEY string
	GOMAIL_EMAIL       string
	GOMAIL_PASSWORD    string
	GOMAIL_HOST        string
	GOMAIL_PORT        int
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

	if val, found := os.LookupEnv("OTP"); found {
		OTP = val
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
			log.Println("error while reading gomail port")
		}
		isRead = false
	}

	if val, found := os.LookupEnv("MIDTRANS_SERVERKEY"); found {
		MIDTRANS_SERVERKEY = val
		isRead = false
	}

	if val, found := os.LookupEnv("GOMAIL_EMAIL"); found {
		GOMAIL_EMAIL = val
		isRead = false
	}

	if val, found := os.LookupEnv("GOMAIL_PASSWORD"); found {
		GOMAIL_PASSWORD = val
		isRead = false
	}

	if val, found := os.LookupEnv("GOMAIL_HOST"); found {
		GOMAIL_HOST = val
		isRead = false
	}

	if val, found := os.LookupEnv("GOMAIL_PORT"); found {
		GOMAIL_PORT, err = strconv.Atoi(val)
		if err != nil {
			log.Println("error while reading gomail port")
		}
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
		OTP = viper.GetString("OTP")
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
		GOMAIL_EMAIL = viper.GetString("GOMAIL_EMAIL")
		GOMAIL_PASSWORD = viper.GetString("GOMAIL_PASSWORD")
		GOMAIL_HOST = viper.GetString("GOMAIL_HOST")
		GOMAIL_PORT = viper.GetInt("GOMAIL_PORT")
	}

	return &app
}
