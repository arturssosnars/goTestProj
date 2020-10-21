package Config

type PostgresConfig struct {
	Database string
	Url string
	Port string
	Driver string
}

type HttpListening struct {
	Port string
}

type BankApi struct {
	Url string
}
