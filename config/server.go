package config

// ServerConfig contains the server configuration
type ServerConfig struct {
	Address string
	Port    string
}

func newServerConfig() (*ServerConfig, error) {
	sc := &ServerConfig{
		Address: get("SERVER_ADDR"),
		Port:    get("PORT"),
	}

	return sc, nil
}
