package config

// Kāpēc tu jūti nepieciešamību datus definēt vienā failā, bet ar to saistīto funkcionalitāti citā failā?
// Es šos struktus definētu blakus metodēm, kas atgriež tos aizpildītus.
// Un kas ir modules?

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
