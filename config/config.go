package config

import (
	"fmt"
	"log"
	"sync"
)

var (
	once      sync.Once
	config    *Config
	sc        *ServerConfig
	dbc       *DBConfig
	configErr error
)

type Config struct {
	Server *ServerConfig
	Db     *DBConfig
}

func Get() *Config {
	return config
}

// Parse parses the config from the environment/preset variables.
func Parse() (*Config, error) {
	once.Do(func() {
		sc, configErr = newServerConfig()
		if configErr != nil {
			return
		}
		dbc, configErr = newDbConfig()
		if configErr != nil {
			return
		}
		config = &Config{
			Server: sc,
			Db:     dbc,
		}
	})
	return config, configErr
}

// get provides a single point to configure config keys:values with defaults.
func get(key string) string {
	switch key {
	case "PG_ADDR":
		return GetEnvWithFallback("PG_ADDR", "127.0.0.1:5432")
	case "PG_USER":
		return GetEnvWithFallback("PG_USER", "postgres")
	case "PG_PASSWORD":
		return GetEnvWithFallback("PG_PASSWORD", "docker")
	case "PG_DATABASE":
		return GetEnvWithFallback("PG_DATABASE", "technoStore")
	case "SERVER_ADDR":
		return GetEnvWithFallback("SERVER_ADDR", "localhost")
	case "PORT":
		return GetEnvWithFallback("PORT", "8080")
	case "DB_LOGGING":
		return GetEnvWithFallback("DB_LOGGING", "false")
	case "DB_MIGRATE":
		return GetEnvWithFallback("DB_MIGRATE", "false")
	}
	log.Fatalf("Undefined config key: %s", key)
	return ""
}

func Print() {
	fmt.Printf(" - %s:                     %s\n", "PG_ADDR", get("PG_ADDR"))
	fmt.Printf(" - %s:                     %s\n", "PG_USER", get("PG_USER"))
	fmt.Printf(" - %s:                 %s\n", "PG_DATABASE", get("PG_DATABASE"))
	fmt.Printf(" - %s:                 %s\n", "SERVER_ADDR", get("SERVER_ADDR"))
	fmt.Printf(" - %s:                 %s\n", "SERVER_PORT", get("PORT"))
	fmt.Printf(" - %s:                  %s\n", "DB_LOGGING", get("DB_LOGGING"))
	fmt.Printf(" - %s:                  %s\n", "DB_MIGRATE", get("DB_MIGRATE"))
}
