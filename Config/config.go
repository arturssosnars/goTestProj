package config

import (
	zlog "github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

//Init initializes project configuration
func Init() error {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("./Config")
	err := viper.ReadInConfig()
	if err != nil {
		zlog.Error().Err(err).Msg("Failed to read config file from ./Config/config.json")
		return err
	}

	return nil
}

//Postgres returns PostgreSQL configuration
func Postgres() PostgresConfig {
	return PostgresConfig{
		URL:      viper.GetString("database.url"),
		Port:     viper.GetString("database.port"),
		Driver:   viper.GetString("database.driver"),
		Database: viper.GetString("database.name"),
	}
}

//Bank returns bank api url
func Bank() BankAPI {
	return BankAPI{viper.GetString("bankApi")}
}

//API returns port for api to listen through
func API() HTTPListening {
	return HTTPListening{viper.GetString("apiPort")}
}
