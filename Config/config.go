package Config

import (
	"fmt"
	"github.com/spf13/viper"
)

func Init() error {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("./Config")
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	return nil
}

func Postgres() PostgresConfig {
	url := fmt.Sprintf("%v", viper.Get("database.url"))
	database := fmt.Sprintf("%v", viper.Get("database.name"))
	driver := fmt.Sprintf("%v", viper.Get("database.driver"))
	apiPort := fmt.Sprintf("%v", viper.Get("database.port"))
	return PostgresConfig{
		Url:      url,
		Port:     apiPort,
		Driver:   driver,
		Database: database,
	}
}

func Bank() BankApi {
	api := fmt.Sprintf("%v", viper.Get("bankApi"))
	return BankApi{api}
}

func Api() HttpListening {
	port := fmt.Sprintf("%v", viper.Get("apiPort"))
	return HttpListening{port}
}
