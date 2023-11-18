package config

type DBConfig struct {
	Address            string
	User               string
	Password           string
	Database           string
	MaxOpenConnections int
	MaxIdleConnections int
	Logging bool
	Migrate bool
}

func newDbConfig() (*DBConfig, error) {
	dbc := &DBConfig{
		Address:  get("PG_ADDR"),
		User:     get("PG_USER"),
		Password: get("PG_PASSWORD"),
		Database: get("PG_DATABASE"),
	}
	if get("DB_LOGGING") == "true" {
		dbc.Logging = true
	} else {
		dbc.Logging = false
	}

	if get("DB_MIGRATE") == "true" {
		dbc.Migrate = true
	} else {
		dbc.Migrate = false
	}

	return dbc, nil
}
