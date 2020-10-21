package config

import (
	"fmt"
	"github.com/spf13/viper"
)

//Init initializes project configuration
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

//Postgres returns PostgreSQL configuration
func Postgres() PostgresConfig {
	// 1. fmt.Sprintf šeit nav nepieciešams, lai dabūtu string vērtību. Iesaku iepazīties ar viper metodēm, kā dabūt
	// atpakaļ string vērtību.

	// 2. Šie variabļi, manuprāt, ir lieki, var vienkārši assignot uzreiz struct fieldā.
	// Piemēram, PostgresConfig{URL: package.GetSomeString(), ...}
	url := fmt.Sprintf("%v", viper.Get("database.url"))
	database := fmt.Sprintf("%v", viper.Get("database.name"))
	driver := fmt.Sprintf("%v", viper.Get("database.driver"))
	apiPort := fmt.Sprintf("%v", viper.Get("database.port"))
	return PostgresConfig{
		URL:      url,
		Port:     apiPort,
		Driver:   driver,
		Database: database,
	}
}

//Bank returns bank api url
func Bank() BankAPI {
	api := fmt.Sprintf("%v", viper.Get("bankApi"))
	return BankAPI{api}
}

//API returns port for api to listen through
func API() HTTPListening {
	port := fmt.Sprintf("%v", viper.Get("apiPort"))
	return HTTPListening{port}
}
