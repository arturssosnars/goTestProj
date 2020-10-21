package config

//PostgresConfig holds PostgreSQL config
type PostgresConfig struct {
	Database string
	URL string
	Port string
	Driver string
}

//HTTPListening holds API port to listen through
type HTTPListening struct {
	Port string
}

//BankAPI holds bank API url
type BankAPI struct {
	URL string
}
